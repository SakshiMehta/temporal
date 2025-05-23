package frontend

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.temporal.io/api/operatorservice/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/dynamicconfig"
	"go.temporal.io/server/common/headers"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/log/tag"
	"go.temporal.io/server/common/metrics"
	"go.temporal.io/server/common/namespace"
	"go.temporal.io/server/common/rpc"
	"go.temporal.io/server/common/rpc/encryption"
	"go.temporal.io/server/common/rpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// HTTPAPIServer is an HTTP API server that forwards requests to gRPC via the
// gRPC interceptors.
type HTTPAPIServer struct {
	server                        http.Server
	listener                      net.Listener
	logger                        log.Logger
	serveMux                      *runtime.ServeMux
	stopped                       chan struct{}
	allowedHosts                  *dynamicconfig.GlobalCachedTypedValue[*regexp.Regexp]
	matchAdditionalHeaders        map[string]bool
	matchAdditionalHeaderPrefixes []string
}

var defaultForwardedHeaders = []string{
	"Authorization-Extras",
	"X-Forwarded-For",
	http.CanonicalHeaderKey(headers.ClientNameHeaderName),
	http.CanonicalHeaderKey(headers.ClientVersionHeaderName),
}

type httpRemoteAddrContextKey struct{}

var (
	errHTTPGRPCListenerNotTCP     = errors.New("must use TCP for gRPC listener to support HTTP API")
	errHTTPGRPCStreamNotSupported = errors.New("stream not supported")
)

// NewHTTPAPIServer creates an [HTTPAPIServer].
//
// routes registered with additionalRouteRegistrationFuncs take precedence over the auto generated grpc proxy routes.
func NewHTTPAPIServer(
	serviceConfig *Config,
	rpcConfig config.RPC,
	grpcListener net.Listener,
	tlsConfigProvider encryption.TLSConfigProvider,
	handler Handler,
	operatorHandler *OperatorHandlerImpl,
	interceptors []grpc.UnaryServerInterceptor,
	metricsHandler metrics.Handler,
	router *mux.Router,
	namespaceRegistry namespace.Registry,
	logger log.Logger,
) (*HTTPAPIServer, error) {
	// Create a TCP listener the same as the frontend one but with different port
	tcpAddrRef, _ := grpcListener.Addr().(*net.TCPAddr)
	if tcpAddrRef == nil {
		return nil, errHTTPGRPCListenerNotTCP
	}
	tcpAddr := *tcpAddrRef
	tcpAddr.Port = rpcConfig.HTTPPort
	var listener net.Listener
	var err error
	if listener, err = net.ListenTCP("tcp", &tcpAddr); err != nil {
		return nil, fmt.Errorf("failed listening for HTTP API on %v: %w", &tcpAddr, err)
	}
	// Close the listener if anything else in this function fails
	success := false
	defer func() {
		if !success {
			_ = listener.Close()
		}
	}()

	// Wrap the listener in a TLS listener if there is any TLS config
	if tlsConfigProvider != nil {
		if tlsConfig, err := tlsConfigProvider.GetFrontendServerConfig(); err != nil {
			return nil, fmt.Errorf("failed getting TLS config for HTTP API: %w", err)
		} else if tlsConfig != nil {
			listener = tls.NewListener(listener, tlsConfig)
		}
	}

	h := &HTTPAPIServer{
		listener:     listener,
		logger:       logger,
		stopped:      make(chan struct{}),
		allowedHosts: serviceConfig.HTTPAllowedHosts,
	}

	// Build 4 possible marshalers in order based on content type
	opts := []runtime.ServeMuxOption{
		runtime.WithMarshalerOption(newTemporalProtoMarshaler("  ", false)),
		runtime.WithMarshalerOption(newTemporalProtoMarshaler("", false)),
		runtime.WithMarshalerOption(newTemporalProtoMarshaler("  ", true)),
		runtime.WithMarshalerOption(newTemporalProtoMarshaler("", true)),
	}

	// Set Temporal service error handler
	opts = append(opts, runtime.WithErrorHandler(h.errorHandler))

	// Match headers w/ default
	h.matchAdditionalHeaders = map[string]bool{}
	for _, v := range defaultForwardedHeaders {
		h.matchAdditionalHeaders[v] = true
	}
	for _, v := range rpcConfig.HTTPAdditionalForwardedHeaders {
		if strings.HasSuffix(v, "*") {
			h.matchAdditionalHeaderPrefixes = append(h.matchAdditionalHeaderPrefixes, http.CanonicalHeaderKey(strings.TrimSuffix(v, "*")))
		} else {
			h.matchAdditionalHeaders[http.CanonicalHeaderKey(v)] = true
		}
	}

	opts = append(opts, runtime.WithMiddlewares(h.allowedHostsMiddleware))
	opts = append(opts, runtime.WithIncomingHeaderMatcher(h.incomingHeaderMatcher))

	// Create inline client connection
	clientConn := newInlineClientConn(
		map[string]any{
			"temporal.api.workflowservice.v1.WorkflowService": handler,
			"temporal.api.operatorservice.v1.OperatorService": operatorHandler,
		},
		interceptors,
		metricsHandler,
		namespaceRegistry,
	)

	// Create serve mux
	h.serveMux = runtime.NewServeMux(opts...)

	err = workflowservice.RegisterWorkflowServiceHandlerClient(
		context.Background(),
		h.serveMux,
		workflowservice.NewWorkflowServiceClient(clientConn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed registering workflowservice HTTP API handler: %w", err)
	}

	err = operatorservice.RegisterOperatorServiceHandlerClient(
		context.Background(),
		h.serveMux,
		operatorservice.NewOperatorServiceClient(clientConn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed registering operatorservice HTTP API handler: %w", err)
	}

	// Set the / handler as our function that wraps serve mux.
	router.PathPrefix("/").HandlerFunc(h.serveHTTP)
	// Register the router as the HTTP server handler.
	h.server.Handler = router

	// Put the remote address on the context
	h.server.ConnContext = func(ctx context.Context, c net.Conn) context.Context {
		return context.WithValue(ctx, httpRemoteAddrContextKey{}, c)
	}

	// We want to set ReadTimeout and WriteTimeout as max idle (and IdleTimeout
	// defaults to ReadTimeout) to ensure that a connection cannot hang over that
	// amount of time.
	h.server.ReadTimeout = serviceConfig.KeepAliveMaxConnectionIdle()
	h.server.WriteTimeout = serviceConfig.KeepAliveMaxConnectionIdle()

	success = true
	return h, nil
}

// Serve serves the HTTP API and does not return until there is a serve error or
// GracefulStop completes. Upon graceful stop, this will return nil. If an error
// is returned, the message is clear that it came from the HTTP API server.
func (h *HTTPAPIServer) Serve() error {
	err := h.server.Serve(h.listener)
	// If the error is for close, we have to wait for the shutdown to complete and
	// we don't consider it an error
	if errors.Is(err, http.ErrServerClosed) {
		<-h.stopped
		err = nil
	}
	// Wrap the error to be clearer it's from the HTTP API
	if err != nil {
		return fmt.Errorf("HTTP API serve failed: %w", err)
	}
	return nil
}

// GracefulStop stops the HTTP server. This will first attempt a graceful stop
// with a drain time, then will hard-stop. This will not return until stopped.
func (h *HTTPAPIServer) GracefulStop(gracefulDrainTime time.Duration) {
	// We try a graceful stop for the amount of time we can drain, then we do a
	// hard stop
	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulDrainTime)
	defer cancel()
	// We intentionally ignore this error, we're gonna stop at this point no
	// matter what. This closes the listener too.
	_ = h.server.Shutdown(shutdownCtx)
	_ = h.server.Close()
	close(h.stopped)
}

func (h *HTTPAPIServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	// Limit the request body to max gRPC size. This is hardcoded to 4MB at the
	// moment using gRPC's default at
	// https://github.com/grpc/grpc-go/blob/0673105ebcb956e8bf50b96e28209ab7845a65ad/server.go#L58
	// which is what the constant is set as at the time of this comment.
	r.Body = http.MaxBytesReader(w, r.Body, rpc.MaxHTTPAPIRequestBytes)

	h.logger.Debug(
		"HTTP API call",
		tag.NewStringTag("http-method", r.Method),
		tag.NewAnyTag("http-url", r.URL),
	)

	// Need to change the accept header based on whether pretty and/or
	// noPayloadShorthand are present
	var acceptHeaderSuffix string
	if _, ok := r.URL.Query()["pretty"]; ok {
		acceptHeaderSuffix += "+pretty"
	}
	if _, ok := r.URL.Query()["noPayloadShorthand"]; ok {
		acceptHeaderSuffix += "+no-payload-shorthand"
	}
	if acceptHeaderSuffix != "" {
		r.Header.Set("Accept", "application/json"+acceptHeaderSuffix)
	}

	// Put the TLS info on the peer context
	if r.TLS != nil {
		var addr net.Addr
		if conn, _ := r.Context().Value(httpRemoteAddrContextKey{}).(net.Conn); conn != nil {
			addr = conn.RemoteAddr()
		}
		r = r.WithContext(peer.NewContext(r.Context(), &peer.Peer{
			Addr: addr,
			AuthInfo: credentials.TLSInfo{
				State:          *r.TLS,
				CommonAuthInfo: credentials.CommonAuthInfo{SecurityLevel: credentials.PrivacyAndIntegrity},
			},
		}))
	}

	// Call gRPC gateway mux
	h.serveMux.ServeHTTP(w, r)
}

func (h *HTTPAPIServer) allowedHostsMiddleware(hf runtime.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		allowedHosts := h.allowedHosts.Get()
		if allowedHosts.MatchString(r.Host) {
			hf(w, r, pathParams)
			return
		}
		w.WriteHeader(http.StatusForbidden)
		// PermissionDenied gRPC code is 7.
		_, _ = w.Write([]byte(`{"code": 7, "message": "Host not allowed"}`))
	}
}

func (h *HTTPAPIServer) errorHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	// Convert the error using serviceerror. The result does not conform to Google
	// gRPC status directly (it conforms to gogo gRPC status), but Err() does
	// based on internal code reading. However, Err() uses Google proto Any
	// which our marshaler is not expecting. So instead we are embedding similar
	// logic to runtime.DefaultHTTPProtoErrorHandler in here but with gogo
	// support. We don't implement custom content type marshaler or trailers at
	// this time.

	s := serviceerror.ToStatus(err)
	w.Header().Set("Content-Type", marshaler.ContentType(struct{}{}))

	sProto := s.Proto()
	buf, merr := marshaler.Marshal(sProto)
	if merr != nil {
		h.logger.Warn("Failed to marshal error message", tag.Error(merr))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code": 13, "message": "failed to marshal error message"}`))
		return
	}

	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	_, _ = w.Write(buf)
}

func (h *HTTPAPIServer) incomingHeaderMatcher(headerName string) (string, bool) {
	// Try ours before falling back to default
	if h.matchAdditionalHeaders[headerName] {
		return headerName, true
	}
	for _, prefix := range h.matchAdditionalHeaderPrefixes {
		if strings.HasPrefix(headerName, prefix) {
			return headerName, true
		}
	}
	return runtime.DefaultHeaderMatcher(headerName)
}

// inlineClientConn is a [grpc.ClientConnInterface] implementation that forwards
// requests directly to gRPC via interceptors. This implementation moves all
// outgoing metadata to incoming and takes resulting outgoing metadata and sets
// as header. But which headers to use and TLS peer context and such are
// expected to be handled by the caller.
type inlineClientConn struct {
	methods           map[string]*serviceMethod
	interceptor       grpc.UnaryServerInterceptor
	requestsCounter   metrics.CounterIface
	namespaceRegistry namespace.Registry
}

var _ grpc.ClientConnInterface = (*inlineClientConn)(nil)

type serviceMethod struct {
	info    grpc.UnaryServerInfo
	handler grpc.UnaryHandler
}

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
var protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

func newInlineClientConn(
	servers map[string]any,
	interceptors []grpc.UnaryServerInterceptor,
	metricsHandler metrics.Handler,
	namespaceRegistry namespace.Registry,
) *inlineClientConn {
	// Create the set of methods via reflection. We currently accept the overhead
	// of reflection compared to having to custom generate gateway code.
	methods := map[string]*serviceMethod{}
	for qualifiedServerName, server := range servers {
		serverVal := reflect.ValueOf(server)
		for i := 0; i < serverVal.Type().NumMethod(); i++ {
			reflectMethod := serverVal.Type().Method(i)
			// We intentionally look this up by name to not assume method indexes line
			// up from type to value
			methodVal := serverVal.MethodByName(reflectMethod.Name)
			// We assume the methods we want only accept a context + request and only
			// return a response + error. We also assume the method name matches the
			// RPC name.
			methodType := methodVal.Type()
			validRPCMethod := methodType.Kind() == reflect.Func &&
				methodType.NumIn() == 2 &&
				methodType.NumOut() == 2 &&
				methodType.In(0) == contextType &&
				methodType.In(1).Implements(protoMessageType) &&
				methodType.Out(0).Implements(protoMessageType) &&
				methodType.Out(1) == errorType
			if !validRPCMethod {
				continue
			}
			fullMethod := "/" + qualifiedServerName + "/" + reflectMethod.Name
			methods[fullMethod] = &serviceMethod{
				info: grpc.UnaryServerInfo{Server: server, FullMethod: fullMethod},
				handler: func(ctx context.Context, req interface{}) (interface{}, error) {
					ret := methodVal.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
					err, _ := ret[1].Interface().(error)
					return ret[0].Interface(), err
				},
			}
		}
	}

	return &inlineClientConn{
		methods:           methods,
		interceptor:       chainUnaryServerInterceptors(interceptors),
		requestsCounter:   metrics.HTTPServiceRequests.With(metricsHandler),
		namespaceRegistry: namespaceRegistry,
	}
}

func (i *inlineClientConn) Invoke(
	ctx context.Context,
	method string,
	args any,
	reply any,
	opts ...grpc.CallOption,
) error {
	// Move outgoing metadata to incoming and set new outgoing metadata
	md, _ := metadata.FromOutgoingContext(ctx)
	// Set the client and version headers if not already set
	if len(md[headers.ClientNameHeaderName]) == 0 {
		md.Set(headers.ClientNameHeaderName, headers.ClientNameServerHTTP)
	}
	if len(md[headers.ClientVersionHeaderName]) == 0 {
		md.Set(headers.ClientVersionHeaderName, headers.ServerVersion)
	}
	ctx = metadata.NewIncomingContext(ctx, md)
	outgoingMD := metadata.MD{}
	ctx = metadata.NewOutgoingContext(ctx, outgoingMD)

	// Get the method. Should never fail, but we check anyways
	serviceMethod := i.methods[method]
	if serviceMethod == nil {
		return status.Error(codes.NotFound, "call not found")
	}

	// Add metric
	var namespaceTag metrics.Tag
	if namespaceName := interceptor.MustGetNamespaceName(i.namespaceRegistry, args); namespaceName != "" {
		namespaceTag = metrics.NamespaceTag(namespaceName.String())
	} else {
		namespaceTag = metrics.NamespaceUnknownTag()
	}
	i.requestsCounter.Record(1, metrics.OperationTag(method), namespaceTag)

	// Invoke
	var resp any
	var err error
	if i.interceptor == nil {
		resp, err = serviceMethod.handler(ctx, args)
	} else {
		resp, err = i.interceptor(ctx, args, &serviceMethod.info, serviceMethod.handler)
	}

	// Find the header call option and set response headers. We accept that if
	// somewhere internally the metadata was replaced instead of appended to, this
	// does not work.
	for _, opt := range opts {
		if callOpt, ok := opt.(grpc.HeaderCallOption); ok {
			*callOpt.HeaderAddr = outgoingMD
		}
	}

	// Merge the response proto onto the wanted reply if non-nil
	if respProto, _ := resp.(proto.Message); respProto != nil {
		proto.Merge(reply.(proto.Message), respProto)
	}

	return err
}

func (*inlineClientConn) NewStream(
	context.Context,
	*grpc.StreamDesc,
	string,
	...grpc.CallOption,
) (grpc.ClientStream, error) {
	return nil, errHTTPGRPCStreamNotSupported
}

// Mostly taken from https://github.com/grpc/grpc-go/blob/v1.56.1/server.go#L1124-L1158
// with slight modifications.
func chainUnaryServerInterceptors(interceptors []grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	switch len(interceptors) {
	case 0:
		return nil
	case 1:
		return interceptors[0]
	default:
		return chainUnaryInterceptors(interceptors)
	}
}

func chainUnaryInterceptors(interceptors []grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return interceptors[0](ctx, req, info, getChainUnaryHandler(interceptors, 0, info, handler))
	}
}

func getChainUnaryHandler(
	interceptors []grpc.UnaryServerInterceptor,
	curr int,
	info *grpc.UnaryServerInfo,
	finalHandler grpc.UnaryHandler,
) grpc.UnaryHandler {
	if curr == len(interceptors)-1 {
		return finalHandler
	}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return interceptors[curr+1](ctx, req, info, getChainUnaryHandler(interceptors, curr+1, info, finalHandler))
	}
}
