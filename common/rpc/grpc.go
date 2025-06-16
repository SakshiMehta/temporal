// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package rpc

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"go.temporal.io/api/serviceerror"
	adminservice "go.temporal.io/server/api/adminservice/v1"
	"go.temporal.io/server/common/headers"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/log/tag"
	"go.temporal.io/server/common/metrics"
	"go.temporal.io/server/common/persistence/serialization"
	"go.temporal.io/server/common/rpc/interceptor"
	serviceerrors "go.temporal.io/server/common/serviceerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// DefaultServiceConfig is a default gRPC connection service config which enables DNS round robin between IPs.
	// To use DNS resolver, a "dns:///" prefix should be applied to the hostPort.
	// https://github.com/grpc/grpc/blob/master/doc/naming.md
	DefaultServiceConfig = `{"loadBalancingConfig": [{"round_robin":{}}]}`

	// MaxBackoffDelay is a maximum interval between reconnect attempts.
	MaxBackoffDelay = 10 * time.Second

	// MaxHTTPAPIRequestBytes is the maximum number of bytes an HTTP API request
	// can have. This is currently set to the max gRPC request size.
	MaxHTTPAPIRequestBytes = 4 * 1024 * 1024

	// MaxNexusAPIRequestBodyBytes is the maximum number of bytes a Nexus HTTP API request can have. Because the body is
	// read into a Payload object, this is currently set to the max Payload size. Content headers are transformed to
	// Payload metadata and contribute to the Payload size as well. A separate limit is enforced on top of this.
	MaxNexusAPIRequestBodyBytes = 2 * 1024 * 1024

	// minConnectTimeout is the minimum amount of time we are willing to give a connection to complete.
	minConnectTimeout = 20 * time.Second

	// maxInternodeRecvPayloadSize indicates the internode max receive payload size.
	maxInternodeRecvPayloadSize = 128 * 1024 * 1024 // 128 Mb

	// ResourceExhaustedCauseHeader will be added to rpc response if request returns ResourceExhausted error.
	// Value of this header will be ResourceExhaustedCause.
	ResourceExhaustedCauseHeader = "X-Resource-Exhausted-Cause"

	// ResourceExhaustedScopeHeader will be added to rpc response if request returns ResourceExhausted error.
	// Value of this header will be the scope of exhausted resource.
	ResourceExhaustedScopeHeader = "X-Resource-Exhausted-Scope"
)

// Dial creates a client connection to the given target with default options.
// The hostName syntax is defined in
// https://github.com/grpc/grpc/blob/master/doc/naming.md.
// dns resolver is used by default
func Dial(hostName string, tlsConfig *tls.Config, logger log.Logger, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	logger.Info("[gRPC] Dial: starting connection", tag.Address(hostName))
	start := time.Now()

	var grpcSecureOpt grpc.DialOption
	var dialOptions []grpc.DialOption

	// Log TLS configuration details
	if tlsConfig == nil {
		logger.Info("[gRPC] Dial: TLS configuration is nil - using insecure credentials",
			tag.Address(hostName),
			tag.NewStringTag("tls_config", "nil"))
	} else {
		logger.Info("[gRPC] Dial: TLS configuration details",
			tag.Address(hostName),
			tag.NewStringTag("server_name", tlsConfig.ServerName),
			tag.NewStringTag("min_version", fmt.Sprintf("%d", tlsConfig.MinVersion)),
			tag.NewStringTag("max_version", fmt.Sprintf("%d", tlsConfig.MaxVersion)),
			tag.NewStringTag("insecure_skip_verify", fmt.Sprintf("%v", tlsConfig.InsecureSkipVerify)),
			tag.NewStringTag("cert_count", fmt.Sprintf("%d", len(tlsConfig.Certificates))),
			tag.NewStringTag("root_ca_count", fmt.Sprintf("%d", len(tlsConfig.RootCAs.Subjects()))),
			tag.NewStringTag("client_ca_count", fmt.Sprintf("%d", len(tlsConfig.ClientCAs.Subjects()))),
			tag.NewStringTag("cipher_suites", fmt.Sprintf("%v", tlsConfig.CipherSuites)),
			tag.NewStringTag("prefer_server_cipher_suites", fmt.Sprintf("%v", tlsConfig.PreferServerCipherSuites)))
	}

	// Handle passthrough addresses specially
	if u, err := url.Parse(hostName); err == nil && u.Scheme == "passthrough" {
		logger.Info("[gRPC] Dial: detected passthrough scheme", tag.Address(hostName))
		hostName = "passthrough:" + strings.TrimPrefix(u.Path, "/")
		logger.Info("[gRPC] Dial: normalized passthrough address", tag.Address(hostName))
		customDialer := func(ctx context.Context, addr string) (net.Conn, error) {
			logger.Info("[gRPC] Dial: custom passthrough dialer invoked", tag.Address(addr))
			dialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			addrWithoutPrefix := strings.TrimPrefix(addr, "passthrough:")
			logger.Info("[gRPC] Dial: custom dialer connecting to", tag.Address(addrWithoutPrefix))
			return dialer.DialContext(ctx, "tcp", addrWithoutPrefix)
		}
		dialOptions = append(dialOptions, grpc.WithContextDialer(customDialer))
	}

	if tlsConfig == nil {
		logger.Info("[gRPC] Dial: using insecure credentials",
			tag.Address(hostName),
			tag.NewStringTag("reason", "tls_config_is_nil"))
		grpcSecureOpt = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		logger.Info("[gRPC] Dial: using TLS credentials",
			tag.Address(hostName),
			tag.NewStringTag("server_name", tlsConfig.ServerName),
			tag.NewStringTag("cert_count", fmt.Sprintf("%d", len(tlsConfig.Certificates))))
		tlsConfigCopy := tlsConfig.Clone()
		host, _, err := net.SplitHostPort(hostName)
		if err != nil {
			host = hostName
		}
		tlsConfigCopy.ServerName = host
		logger.Info("[gRPC] Dial: setting TLS server name",
			tag.Address(hostName),
			tag.NewStringTag("original_server_name", tlsConfig.ServerName),
			tag.NewStringTag("new_server_name", host))
		grpcSecureOpt = grpc.WithTransportCredentials(credentials.NewTLS(tlsConfigCopy))
	}

	var cp = grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: minConnectTimeout,
	}
	cp.Backoff.MaxDelay = MaxBackoffDelay

	dialOptions = append(dialOptions,
		grpcSecureOpt,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxInternodeRecvPayloadSize)),
		grpc.WithChainUnaryInterceptor(
			headersInterceptor,
			metrics.NewClientMetricsTrailerPropagatorInterceptor(logger),
			errorInterceptor,
		),
		grpc.WithChainStreamInterceptor(
			interceptor.StreamErrorInterceptor,
		),
		grpc.WithDefaultServiceConfig(DefaultServiceConfig),
		grpc.WithDisableServiceConfig(),
		grpc.WithConnectParams(cp),
	)
	dialOptions = append(dialOptions, opts...)

	// Add our custom logging interceptor
	dialOptions = append(dialOptions, grpc.WithUnaryInterceptor(
		adminServiceLoggingInterceptor(logger),
	))

	logger.Info("[gRPC] Dial: dialing grpc.NewClient", tag.Address(hostName))
	conn, err := grpc.NewClient(hostName, dialOptions...)
	if err != nil {
		logger.Error("[gRPC] Dial: failed to create gRPC connection", tag.Error(err), tag.Address(hostName))
		return nil, err
	}
	logger.Info("[gRPC] Dial: successfully created gRPC connection", tag.Address(hostName), tag.NewStringTag("duration", time.Since(start).String()))
	return conn, nil
}

func errorInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	err = serviceerrors.FromStatus(status.Convert(err))
	return err
}

func headersInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	ctx = headers.Propagate(ctx)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func ServiceErrorInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	resp, err := handler(ctx, req)

	var deserializationError *serialization.DeserializationError
	var serializationError *serialization.SerializationError
	// convert serialization errors to be captured as serviceerrors across gRPC calls
	if errors.As(err, &deserializationError) || errors.As(err, &serializationError) {
		err = serviceerror.NewDataLoss(err.Error())
	}
	return resp, serviceerror.ToStatus(err).Err()
}

func NewFrontendServiceErrorInterceptor(
	logger log.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		resp, err := handler(ctx, req)

		if err == nil {
			return resp, err
		}

		// mask some internal service errors at frontend
		switch err.(type) {
		case *serviceerrors.ShardOwnershipLost:
			err = serviceerror.NewUnavailable("shard unavailable, please backoff and retry")
		case *serviceerror.DataLoss:
			err = serviceerror.NewUnavailable("internal history service error")
		}

		addHeadersForResourceExhausted(ctx, logger, err)

		return resp, err
	}
}

func addHeadersForResourceExhausted(ctx context.Context, logger log.Logger, err error) {
	var reErr *serviceerror.ResourceExhausted
	if errors.As(err, &reErr) {
		headerErr := grpc.SetHeader(ctx, metadata.Pairs(
			ResourceExhaustedCauseHeader, reErr.Cause.String(),
			ResourceExhaustedScopeHeader, reErr.Scope.String(),
		))
		if headerErr != nil {
			logger.Error("Failed to add Resource-Exhausted headers to response", tag.Error(headerErr))
		}
	}
}

// adminServiceLoggingInterceptor adds detailed logging for admin service RPC calls
func adminServiceLoggingInterceptor(logger log.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := time.Now()
		connState := cc.GetState().String()

		// Log before the call
		logger.Info("[AdminService] Starting RPC call",
			tag.NewStringTag("method", method),
			tag.NewStringTag("target", cc.Target()),
			tag.NewStringTag("request_type", fmt.Sprintf("%T", req)),
			tag.NewStringTag("connection_state", connState))

		// For DescribeCluster specifically, log detailed request information
		if method == "/temporal.server.api.adminservice.v1.AdminService/DescribeCluster" {
			if req, ok := req.(*adminservice.DescribeClusterRequest); ok {
				logger.Info("[AdminService] DescribeCluster request details",
					tag.NewStringTag("cluster_name", req.GetClusterName()),
					tag.NewStringTag("target", cc.Target()),
					tag.NewStringTag("connection_state", connState),
					tag.NewStringTag("request_id", fmt.Sprintf("%p", req)), // Add request ID for correlation
					tag.NewStringTag("timestamp", startTime.Format(time.RFC3339)))
			}
		}

		// Make the actual call
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(startTime)

		// Log after the call
		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("[AdminService] RPC call failed",
				tag.Error(err),
				tag.NewStringTag("method", method),
				tag.NewStringTag("target", cc.Target()),
				tag.NewStringTag("error_type", fmt.Sprintf("%T", err)),
				tag.NewStringTag("error_code", st.Code().String()),
				tag.NewStringTag("duration", duration.String()),
				tag.NewStringTag("connection_state", cc.GetState().String()))

			// For DescribeCluster specifically, log more error details
			if method == "/temporal.server.api.adminservice.v1.AdminService/DescribeCluster" {
				logger.Error("[AdminService] DescribeCluster call failed",
					tag.Error(err),
					tag.NewStringTag("target", cc.Target()),
					tag.NewStringTag("error_type", fmt.Sprintf("%T", err)),
					tag.NewStringTag("error_code", st.Code().String()),
					tag.NewStringTag("error_details", err.Error()),
					tag.NewStringTag("duration", duration.String()),
					tag.NewStringTag("connection_state", cc.GetState().String()),
					tag.NewStringTag("request_id", fmt.Sprintf("%p", req))) // Add request ID for correlation
			}
		} else {
			logger.Info("[AdminService] RPC call succeeded",
				tag.NewStringTag("method", method),
				tag.NewStringTag("target", cc.Target()),
				tag.NewStringTag("duration", duration.String()),
				tag.NewStringTag("connection_state", cc.GetState().String()))

			// For DescribeCluster specifically, log detailed response information
			if method == "/temporal.server.api.adminservice.v1.AdminService/DescribeCluster" {
				if resp, ok := reply.(*adminservice.DescribeClusterResponse); ok {
					logger.Info("[AdminService] DescribeCluster response details",
						tag.NewStringTag("cluster_name", resp.GetClusterName()),
						tag.NewStringTag("cluster_id", resp.GetClusterId()),
						tag.NewStringTag("target", cc.Target()),
						tag.NewStringTag("duration", duration.String()),
						tag.NewStringTag("connection_state", cc.GetState().String()),
						tag.NewStringTag("history_shard_count", fmt.Sprintf("%d", resp.GetHistoryShardCount())),
						tag.NewStringTag("persistence_store", resp.GetPersistenceStore()),
						tag.NewStringTag("visibility_store", resp.GetVisibilityStore()),
						tag.NewStringTag("request_id", fmt.Sprintf("%p", req))) // Add request ID for correlation
				}
			}
		}

		return err
	}
}
