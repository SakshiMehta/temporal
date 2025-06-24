// The MIT License
//
// Copyright (c) 2024 Temporal Technologies Inc.  All rights reserved.
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

package cluster

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"go.temporal.io/api/serviceerror"
	"go.temporal.io/server/common"
	"go.temporal.io/server/common/collection"
	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/log/tag"
)

type tlsConfigProvider interface {
	GetRemoteClusterClientConfig(hostname string) (*tls.Config, error)
}

type FrontendHTTPClientCache struct {
	metadata    Metadata
	tlsProvider tlsConfigProvider
	clients     *collection.FallibleOnceMap[string, *common.FrontendHTTPClient]
	logger      log.Logger
}

func NewFrontendHTTPClientCache(
	metadata Metadata,
	tlsProvider tlsConfigProvider,
	logger log.Logger,
) *FrontendHTTPClientCache {
	cache := &FrontendHTTPClientCache{
		metadata:    metadata,
		tlsProvider: tlsProvider,
		logger:      logger,
	}
	cache.clients = collection.NewFallibleOnceMap(cache.newClientForCluster)
	metadata.RegisterMetadataChangeCallback(cache, cache.evictionCallback)
	return cache
}

// Get returns a cached HttpClient if available, or constructs a new one for the given cluster name.
func (c *FrontendHTTPClientCache) Get(targetClusterName string) (*common.FrontendHTTPClient, error) {
	return c.clients.Get(targetClusterName)
}

func (c *FrontendHTTPClientCache) newClientForCluster(targetClusterName string) (*common.FrontendHTTPClient, error) {
	c.logger.Info("newClientForCluster called", tag.NewStringTag("targetClusterName", targetClusterName))

	targetInfo, ok := c.metadata.GetAllClusterInfo()[targetClusterName]
	if !ok {
		c.logger.Error("Cluster metadata not found", tag.NewStringTag("targetClusterName", targetClusterName))
		return nil, serviceerror.NewNotFound(fmt.Sprintf("could not find cluster metadata for cluster %s", targetClusterName))
	}

	c.logger.Info("Found cluster metadata", tag.NewStringTag("targetClusterName", targetClusterName), tag.NewStringTag("httpAddress", targetInfo.HTTPAddress))

	if targetInfo.HTTPAddress == "" {
		c.logger.Error("HTTPAddress not configured for cluster", tag.NewStringTag("targetClusterName", targetClusterName))
		return nil, serviceerror.NewInternal(fmt.Sprintf("HTTPAddress not configured for cluster: %s", targetClusterName))
	}
	host, port, err := net.SplitHostPort(targetInfo.HTTPAddress)
	if err != nil {
		c.logger.Error("Invalid frontend address for cluster", tag.NewStringTag("targetClusterName", targetClusterName), tag.NewStringTag("httpAddress", targetInfo.HTTPAddress), tag.Error(err))
		return nil, fmt.Errorf("%w: %w", serviceerror.NewInternal("invalid frontend address"), err)
	}
	c.logger.Info("SplitHostPort result", tag.NewStringTag("host", host), tag.NewStringTag("port", port))

	client := http.Client{}

	urlScheme := "http"
	if c.tlsProvider != nil {
		c.logger.Info("tlsProvider is present, attempting to get TLS config", tag.NewStringTag("host", host))
		tlsClientConfig, err := c.tlsProvider.GetRemoteClusterClientConfig(host)
		if err != nil {
			c.logger.Error("Failed to get TLS config from tlsProvider", tag.NewStringTag("host", host), tag.Error(err))
			return nil, err
		}
		client.Transport = &http.Transport{TLSClientConfig: tlsClientConfig}
		if tlsClientConfig != nil {
			urlScheme = "https"
			c.logger.Info("TLS config found, using https scheme", tag.NewStringTag("host", host))
		} else {
			c.logger.Info("No TLS config found, using http scheme", tag.NewStringTag("host", host))
		}
	} else {
		c.logger.Info("No tlsProvider present, using http scheme", tag.NewStringTag("host", host))
	}

	c.logger.Info("Returning FrontendHTTPClient", tag.NewStringTag("address", targetInfo.HTTPAddress), tag.NewStringTag("scheme", urlScheme))
	return &common.FrontendHTTPClient{
		Address: targetInfo.HTTPAddress,
		Scheme:  urlScheme,
		Client:  client,
	}, nil
}

// evictionCallback is invoked by cluster.Metadata when cluster information changes.
// It invalidates clients which are either no longer present or have had their HTTP address changed.
// It is assumed that TLS information has not changed for clusters that are unmodified.
func (c *FrontendHTTPClientCache) evictionCallback(oldClusterMetadata map[string]*ClusterInformation, newClusterMetadata map[string]*ClusterInformation) {
	for oldClusterName, oldClusterInfo := range oldClusterMetadata {
		if oldClusterName == c.metadata.GetCurrentClusterName() || oldClusterInfo == nil {
			continue
		}

		newClusterInfo, exists := newClusterMetadata[oldClusterName]
		if !exists || newClusterInfo == nil || oldClusterInfo.HTTPAddress != newClusterInfo.HTTPAddress {
			// Cluster was removed or had its HTTP address changed, so invalidate the cached client for that cluster.
			client, ok := c.clients.Pop(oldClusterName)
			if ok {
				client.CloseIdleConnections()
			}
		}
	}
}
