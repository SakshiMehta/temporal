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
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"go.temporal.io/api/serviceerror"
	"go.temporal.io/server/common"
	"go.temporal.io/server/common/cache"
	"go.temporal.io/server/common/config"
	"go.temporal.io/server/common/convert"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/log/tag"
	"go.temporal.io/server/common/membership"
	"go.temporal.io/server/common/primitives"
	"go.temporal.io/server/common/rpc/encryption"
	"go.temporal.io/server/environment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var _ common.RPCFactory = (*RPCFactory)(nil)

// RPCFactory is an implementation of common.RPCFactory interface
type RPCFactory struct {
	config      *config.RPC
	serviceName primitives.ServiceName
	logger      log.Logger

	frontendURL       string
	frontendHTTPURL   string
	frontendHTTPPort  int
	frontendTLSConfig *tls.Config

	grpcListener func() net.Listener
	tlsFactory   encryption.TLSConfigProvider
	dialOptions  []grpc.DialOption
	monitor      membership.Monitor
	// A OnceValues wrapper for createLocalFrontendHTTPClient.
	localFrontendClient      func() (*common.FrontendHTTPClient, error)
	interNodeGrpcConnections cache.Cache
}

// NewFactory builds a new RPCFactory
// conforming to the underlying configuration
func NewFactory(
	cfg *config.RPC,
	sName primitives.ServiceName,
	logger log.Logger,
	tlsProvider encryption.TLSConfigProvider,
	frontendURL string,
	frontendHTTPURL string,
	frontendHTTPPort int,
	frontendTLSConfig *tls.Config,
	dialOptions []grpc.DialOption,
	monitor membership.Monitor,
) *RPCFactory {
	f := &RPCFactory{
		config:            cfg,
		serviceName:       sName,
		logger:            logger,
		frontendURL:       frontendURL,
		frontendHTTPURL:   frontendHTTPURL,
		frontendHTTPPort:  frontendHTTPPort,
		frontendTLSConfig: frontendTLSConfig,
		tlsFactory:        tlsProvider,
		dialOptions:       dialOptions,
		monitor:           monitor,
	}
	f.grpcListener = sync.OnceValue(f.createGRPCListener)
	f.localFrontendClient = sync.OnceValues(f.createLocalFrontendHTTPClient)
	f.interNodeGrpcConnections = cache.NewSimple(nil)
	return f
}

func (d *RPCFactory) GetFrontendGRPCServerOptions() ([]grpc.ServerOption, error) {
	var opts []grpc.ServerOption

	if d.tlsFactory != nil {
		serverConfig, err := d.tlsFactory.GetFrontendServerConfig()
		if err != nil {
			return nil, err
		}
		if serverConfig == nil {
			return opts, nil
		}
		opts = append(opts, grpc.Creds(credentials.NewTLS(serverConfig)))
	}

	return opts, nil
}

func (d *RPCFactory) GetFrontendClientTlsConfig() (*tls.Config, error) {
	if d.tlsFactory != nil {
		return d.tlsFactory.GetFrontendClientConfig()
	}

	return nil, nil
}

func (d *RPCFactory) GetRemoteClusterClientConfig(hostname string) (*tls.Config, error) {
	d.logger.Info("GetRemoteClusterClientConfig called in RPCFactory", tag.NewStringTag("hostname", hostname))

	if d.tlsFactory != nil {
		d.logger.Info("tlsFactory is present, delegating to tlsFactory.GetRemoteClusterClientConfig")
		config, err := d.tlsFactory.GetRemoteClusterClientConfig(hostname)
		if err != nil {
			d.logger.Error("Error from tlsFactory.GetRemoteClusterClientConfig", tag.Error(err), tag.NewStringTag("hostname", hostname))
			return nil, err
		}
		if config != nil {
			d.logger.Info("tlsFactory returned a TLS config", tag.NewStringTag("serverName", config.ServerName))
		} else {
			d.logger.Info("tlsFactory returned nil TLS config")
		}
		return config, nil
	}

	d.logger.Info("No tlsFactory present, returning nil TLS config")
	return nil, nil
}

func (d *RPCFactory) GetInternodeGRPCServerOptions() ([]grpc.ServerOption, error) {
	var opts []grpc.ServerOption

	if d.tlsFactory != nil {
		serverConfig, err := d.tlsFactory.GetInternodeServerConfig()
		if err != nil {
			return nil, err
		}
		if serverConfig == nil {
			return opts, nil
		}
		opts = append(opts, grpc.Creds(credentials.NewTLS(serverConfig)))
	}

	return opts, nil
}

func (d *RPCFactory) GetInternodeClientTlsConfig() (*tls.Config, error) {
	if d.tlsFactory != nil {
		return d.tlsFactory.GetInternodeClientConfig()
	}

	return nil, nil
}

// GetGRPCListener returns cached dispatcher for gRPC inbound or creates one
func (d *RPCFactory) GetGRPCListener() net.Listener {
	return d.grpcListener()
}

func (d *RPCFactory) createGRPCListener() net.Listener {
	hostAddress := net.JoinHostPort(getListenIP(d.config, d.logger).String(), convert.IntToString(d.config.GRPCPort))

	grpcListener, err := net.Listen("tcp", hostAddress)
	if err != nil || grpcListener == nil || grpcListener.Addr() == nil {
		d.logger.Fatal("Failed to start gRPC listener", tag.Error(err), tag.Service(d.serviceName), tag.Address(hostAddress))
	}

	d.logger.Info("Created gRPC listener", tag.Service(d.serviceName), tag.Address(hostAddress))
	return grpcListener
}

func getListenIP(cfg *config.RPC, logger log.Logger) net.IP {
	if cfg.BindOnLocalHost && len(cfg.BindOnIP) > 0 {
		logger.Fatal("ListenIP failed, bindOnLocalHost and bindOnIP are mutually exclusive")
		return nil
	}

	if cfg.BindOnLocalHost {
		return net.ParseIP(environment.GetLocalhostIP())
	}

	if len(cfg.BindOnIP) > 0 {
		ip := net.ParseIP(cfg.BindOnIP)
		if ip != nil {
			return ip
		}
		logger.Fatal("ListenIP failed, unable to parse bindOnIP value", tag.Address(cfg.BindOnIP))
		return nil
	}
	ip, err := config.ListenIP()
	if err != nil {
		logger.Fatal("ListenIP failed", tag.Error(err))
		return nil
	}
	return ip
}

// CreateRemoteFrontendGRPCConnection creates connection for gRPC calls
func (d *RPCFactory) CreateRemoteFrontendGRPCConnection(rpcAddress string) *grpc.ClientConn {
	d.logger.Info("Starting CreateRemoteFrontendGRPCConnection",
		tag.Address(rpcAddress),
		tag.NewStringTag("connection_type", "remote_frontend"))

	// Determine the actual address to connect to and the hostname for TLS config
	target := rpcAddress
	hostname := ""
	var err error

	if strings.HasPrefix(rpcAddress, "passthrough:") {
		target = strings.TrimPrefix(rpcAddress, "passthrough:")
		hostname, _, err = net.SplitHostPort(target)
		if err != nil {
			d.logger.Warn("Failed to split host and port after stripping passthrough prefix, using 'passthrough' as hostname", tag.Error(err), tag.Address(target))
			hostname = "passthrough"
		}
	} else {
		hostname, _, err = net.SplitHostPort(target)
		if err != nil {
			d.logger.Fatal("Invalid rpcAddress for remote cluster. Unable to extract host:port", tag.Error(err), tag.Address(target))
		}
	}

	// Prepare TLS config if needed
	var tlsClientConfig *tls.Config
	if d.tlsFactory != nil {
		d.logger.Info("Extracted hostname for TLS config",
			tag.NewStringTag("hostname", hostname),
			tag.Address(target))
		tlsClientConfig, err = d.tlsFactory.GetRemoteClusterClientConfig(hostname)
		if err != nil {
			d.logger.Fatal("Failed to create tls config for gRPC connection", tag.Error(err), tag.NewStringTag("hostname", hostname), tag.Address(target))
			return nil
		}
		if tlsClientConfig != nil {
			d.logger.Info("TLS config loaded successfully",
				tag.NewStringTag("serverName", tlsClientConfig.ServerName),
				tag.NewStringTag("min_version", fmt.Sprintf("%d", tlsClientConfig.MinVersion)),
				tag.NewStringTag("max_version", fmt.Sprintf("%d", tlsClientConfig.MaxVersion)),
				tag.NewBoolTag("insecure_skip_verify", tlsClientConfig.InsecureSkipVerify),
				tag.NewStringTag("cert_count", fmt.Sprintf("%d", len(tlsClientConfig.))),
				tag.NewStringTag("root_ca_count", func() string {
					if tlsClientConfig.RootCAs != nil {
						return fmt.Sprintf("%d", len(tlsClientConfig.RootCAs.Subjects()))
					}
					return "0"
				}()),
				tag.Address(target))
		} else {
			d.logger.Info("No TLS config, using insecure connection", tag.Address(target))
		}
	}

	// Establish the gRPC connection
	conn := d.dial(target, tlsClientConfig)
	if conn != nil {
		d.logger.Info("Successfully created gRPC connection",
			tag.Address(target),
			tag.NewStringTag("connection_state", conn.GetState().String()))
	} else {
		d.logger.Error("Failed to create gRPC connection", tag.Address(target))
	}
	return conn
}

// CreateLocalFrontendGRPCConnection creates connection for internal frontend calls
func (d *RPCFactory) CreateLocalFrontendGRPCConnection() *grpc.ClientConn {
	return d.dial(d.frontendURL, d.frontendTLSConfig)
}

// CreateInternodeGRPCConnection creates connection for gRPC calls
func (d *RPCFactory) CreateInternodeGRPCConnection(hostName string) *grpc.ClientConn {
	if c, ok := d.interNodeGrpcConnections.Get(hostName).(*grpc.ClientConn); ok {
		return c
	}
	var tlsClientConfig *tls.Config
	var err error
	if d.tlsFactory != nil {
		tlsClientConfig, err = d.tlsFactory.GetInternodeClientConfig()
		if err != nil {
			d.logger.Fatal("Failed to create tls config for gRPC connection", tag.Error(err))
			return nil
		}
	}
	c := d.dial(hostName, tlsClientConfig)
	d.interNodeGrpcConnections.Put(hostName, c)
	return c
}

func (d *RPCFactory) dial(hostName string, tlsClientConfig *tls.Config) *grpc.ClientConn {
	connection, err := Dial(hostName, tlsClientConfig, d.logger, d.dialOptions...)
	if err != nil {
		d.logger.Fatal("Failed to create gRPC connection", tag.Error(err))
		return nil
	}

	return connection
}

func (d *RPCFactory) GetTLSConfigProvider() encryption.TLSConfigProvider {
	return d.tlsFactory
}

// CreateLocalFrontendHTTPClient gets or creates a cached frontend client.
func (d *RPCFactory) CreateLocalFrontendHTTPClient() (*common.FrontendHTTPClient, error) {
	return d.localFrontendClient()
}

// createLocalFrontendHTTPClient creates an HTTP client for communicating with the frontend.
// It uses either the provided frontendURL or membership to resolve the frontend address.
func (d *RPCFactory) createLocalFrontendHTTPClient() (*common.FrontendHTTPClient, error) {
	// dialer and transport field values copied from http.DefaultTransport.
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := http.Client{}

	// Default to http unless TLS is configured.
	scheme := "http"
	if d.frontendTLSConfig != nil {
		transport.TLSClientConfig = d.frontendTLSConfig
		scheme = "https"
	}

	var address string
	if r := serviceResolverFromGRPCURL(d.frontendHTTPURL); r != nil {
		client.Transport = &roundTripper{
			resolver:   r,
			underlying: transport,
			httpPort:   d.frontendHTTPPort,
		}
		address = "internal" // This will be replaced by the roundTripper
	} else {
		// Use the URL as-is and leave the transport unmodified.
		client.Transport = transport
		address = d.frontendHTTPURL
	}

	return &common.FrontendHTTPClient{
		Client:  client,
		Address: address,
		Scheme:  scheme,
	}, nil
}

type roundTripper struct {
	resolver   membership.ServiceResolver
	underlying http.RoundTripper
	httpPort   int
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Pick a frontend host at random.
	members := rt.resolver.AvailableMembers()
	if len(members) == 0 {
		return nil, serviceerror.NewUnavailable("no frontend host to route request to")
	}
	idx := rand.Intn(len(members))
	member := members[idx]

	// Replace port with the HTTP port.
	host, _, err := net.SplitHostPort(member.Identity())
	if err != nil {
		return nil, fmt.Errorf("failed to extract port from frontend member: %w", err)
	}
	address := fmt.Sprintf("%s:%d", host, rt.httpPort)

	// Replace request's host.
	req.URL.Host = address
	req.Host = address
	return rt.underlying.RoundTrip(req)
}

// serviceResolverFromGRPCURL returns a ServiceResolver if ustr corresponds to a
// membership url, otherwise nil.
func serviceResolverFromGRPCURL(ustr string) membership.ServiceResolver {
	u, err := url.Parse(ustr)
	if err != nil {
		return nil
	}
	res, err := membership.GetServiceResolverFromURL(u)
	if err != nil {
		return nil
	}
	return res
}
