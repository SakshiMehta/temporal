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
	"go.temporal.io/server/common/headers"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/log/tag"
	"go.temporal.io/server/common/persistence/serialization"
	serviceerrors "go.temporal.io/server/common/serviceerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
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
	logger.Info("Starting gRPC Dial",
		tag.Address(hostName),
		tag.NewStringTag("dial_function", "common/rpc/grpc.Dial"))

	start := time.Now()

	var grpcSecureOpt grpc.DialOption
	var dialOptions []grpc.DialOption

	// Log TLS configuration details
	if tlsConfig == nil {
		logger.Info("TLS configuration is nil - using insecure credentials",
			tag.Address(hostName),
			tag.NewStringTag("tls_config", "nil"))
	} else {
		logger.Info("TLS configuration details",
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

	// Add DNS resolution logging
	if host, port, err := net.SplitHostPort(hostName); err == nil {
		logger.Info("Attempting DNS resolution",
			tag.NewStringTag("host", host),
			tag.NewStringTag("port", port),
			tag.Address(hostName))

		if ips, err := net.LookupHost(host); err == nil {
			logger.Info("DNS resolution successful",
				tag.NewStringTag("host", host),
				tag.NewStringTag("resolved_ips", strings.Join(ips, ",")),
				tag.Address(hostName))

			// Log each resolved IP with more detail
			for i, ip := range ips {
				logger.Info("Resolved IP details",
					tag.NewStringTag("ip_index", fmt.Sprintf("%d", i)),
					tag.NewStringTag("ip", ip),
					tag.NewStringTag("host", host),
					tag.Address(hostName))

				// Try to get more info about the IP
				if parsedIP := net.ParseIP(ip); parsedIP != nil {
					logger.Info("IP characteristics",
						tag.NewStringTag("ip", ip),
						tag.NewStringTag("is_loopback", fmt.Sprintf("%t", parsedIP.IsLoopback())),
						tag.NewStringTag("is_private", fmt.Sprintf("%t", parsedIP.IsPrivate())),
						tag.NewStringTag("is_global_unicast", fmt.Sprintf("%t", parsedIP.IsGlobalUnicast())),
						tag.Address(hostName))
				}
			}
		} else {
			logger.Error("DNS resolution failed",
				tag.Error(err),
				tag.NewStringTag("host", host),
				tag.Address(hostName))
		}
	} else {
		logger.Warn("Could not split host:port for DNS resolution",
			tag.Error(err),
			tag.Address(hostName))
	}

	// Handle passthrough addresses specially
	if u, err := url.Parse(hostName); err == nil && u.Scheme == "passthrough" {
		logger.Info("Detected passthrough scheme", tag.Address(hostName))
		hostName = "passthrough:" + strings.TrimPrefix(u.Path, "/")
		logger.Info("Normalized passthrough address", tag.Address(hostName))
		customDialer := func(ctx context.Context, addr string) (net.Conn, error) {
			logger.Info("Custom passthrough dialer invoked", tag.Address(addr))
			dialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}
			addrWithoutPrefix := strings.TrimPrefix(addr, "passthrough:")
			logger.Info("Custom dialer connecting to", tag.Address(addrWithoutPrefix))
			return dialer.DialContext(ctx, "tcp", addrWithoutPrefix)
		}
		dialOptions = append(dialOptions, grpc.WithContextDialer(customDialer))
	}

	if tlsConfig == nil {
		logger.Info("Using insecure credentials",
			tag.Address(hostName),
			tag.NewStringTag("reason", "tls_config_is_nil"))
		grpcSecureOpt = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		logger.Info("Using TLS credentials",
			tag.Address(hostName),
			tag.NewStringTag("server_name", tlsConfig.ServerName),
			tag.NewStringTag("cert_count", fmt.Sprintf("%d", len(tlsConfig.Certificates))))
		tlsConfigCopy := tlsConfig.Clone()
		host, _, err := net.SplitHostPort(hostName)
		if err != nil {
			host = hostName
		}
		tlsConfigCopy.ServerName = host
		logger.Info("Setting TLS server name",
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
			errorInterceptor,
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

	logger.Info("Dialing grpc.NewClient",
		tag.Address(hostName),
		tag.NewStringTag("dial_options_count", fmt.Sprintf("%d", len(dialOptions))))

	conn, err := grpc.NewClient(hostName, dialOptions...)
	if err != nil {
		logger.Error("Failed to create gRPC connection",
			tag.Error(err),
			tag.Address(hostName),
			tag.NewStringTag("error_type", fmt.Sprintf("%T", err)),
			tag.NewStringTag("duration", time.Since(start).String()))
		return nil, err
	}

	logger.Info("Successfully created gRPC connection",
		tag.Address(hostName),
		tag.NewStringTag("duration", time.Since(start).String()),
		tag.NewStringTag("connection_state", conn.GetState().String()))

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
		connTarget := cc.Target()

		// Log connection state before the call
		logger.Info("[AdminService] Connection state before RPC call",
			tag.NewStringTag("method", method),
			tag.NewStringTag("target", connTarget),
			tag.NewStringTag("connection_state", connState),
			tag.NewStringTag("connection_id", fmt.Sprintf("%p", cc)))

		// For PollActivityTaskQueue specifically, add more detailed logging
		if method == "/temporal.server.api.matchingservice.v1.MatchingService/PollActivityTaskQueue" {
			logger.Info("[AdminService] PollActivityTaskQueue request details",
				tag.NewStringTag("target", connTarget),
				tag.NewStringTag("connection_state", connState),
				tag.NewStringTag("connection_id", fmt.Sprintf("%p", req)))
		}

		// Make the actual call
		err := invoker(ctx, method, req, reply, cc, opts...)
		duration := time.Since(startTime)
		newConnState := cc.GetState().String()

		// Log connection state after the call
		if newConnState != connState {
			logger.Info("[AdminService] Connection state changed during RPC call",
				tag.NewStringTag("method", method),
				tag.NewStringTag("target", connTarget),
				tag.NewStringTag("old_state", connState),
				tag.NewStringTag("new_state", newConnState),
				tag.NewStringTag("connection_id", fmt.Sprintf("%p", cc)),
				tag.NewStringTag("duration", duration.String()))
		}

		// Log after the call
		if err != nil {
			st, _ := status.FromError(err)
			// Check for specific error types
			if st.Code() == codes.Unavailable {
				logger.Error("[AdminService] RPC call failed - service unavailable",
					tag.Error(err),
					tag.NewStringTag("method", method),
					tag.NewStringTag("target", connTarget),
					tag.NewStringTag("error_type", fmt.Sprintf("%T", err)),
					tag.NewStringTag("error_code", st.Code().String()),
					tag.NewStringTag("duration", duration.String()),
					tag.NewStringTag("connection_state", newConnState),
					tag.NewStringTag("connection_id", fmt.Sprintf("%p", cc)),
					tag.NewStringTag("error_details", st.Message()),
					tag.NewStringTag("debug_data", fmt.Sprintf("%v", st.Details())))

				// For PollActivityTaskQueue, add specific graceful shutdown detection
				if method == "/temporal.server.api.matchingservice.v1.MatchingService/PollActivityTaskQueue" {
					if strings.Contains(st.Message(), "graceful_stop") {
						logger.Error("[AdminService] PollActivityTaskQueue failed due to graceful shutdown",
							tag.NewStringTag("target", connTarget),
							tag.NewStringTag("connection_state", newConnState),
							tag.NewStringTag("connection_id", fmt.Sprintf("%p", cc)),
							tag.NewStringTag("duration", duration.String()),
							tag.NewStringTag("error_details", st.Message()))
					}
				}
			} else {
				logger.Error("[AdminService] RPC call failed",
					tag.Error(err),
					tag.NewStringTag("method", method),
					tag.NewStringTag("target", connTarget),
					tag.NewStringTag("error_type", fmt.Sprintf("%T", err)),
					tag.NewStringTag("error_code", st.Code().String()),
					tag.NewStringTag("duration", duration.String()),
					tag.NewStringTag("connection_state", newConnState),
					tag.NewStringTag("connection_id", fmt.Sprintf("%p", cc)))
			}
		} else {
			logger.Info("[AdminService] RPC call succeeded",
				tag.NewStringTag("method", method),
				tag.NewStringTag("target", connTarget),
				tag.NewStringTag("duration", duration.String()),
				tag.NewStringTag("connection_state", newConnState),
				tag.NewStringTag("connection_id", fmt.Sprintf("%p", cc)))
		}

		return err
	}
}
