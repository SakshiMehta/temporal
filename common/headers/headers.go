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

package headers

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

// Note the nexusoperations component references these headers and adds them to a list of disallowed headers for users to set.
// If any other headers are added for internal use, they should be added to the disallowed headers list.
const (
	ClientNameHeaderName              = "client-name"
	ClientVersionHeaderName           = "client-version"
	SupportedServerVersionsHeaderName = "supported-server-versions"
	SupportedFeaturesHeaderName       = "supported-features"
	SupportedFeaturesHeaderDelim      = ","

	CallerNameHeaderName = "caller-name"
	CallerTypeHeaderName = "caller-type"
	CallOriginHeaderName = "call-initiation"
)

var (
	// propagateHeaders are the headers to propagate from the frontend to other services.
	propagateHeaders = []string{
		ClientNameHeaderName,
		ClientVersionHeaderName,
		SupportedServerVersionsHeaderName,
		SupportedFeaturesHeaderName,
		CallerNameHeaderName,
		CallerTypeHeaderName,
		CallOriginHeaderName,
		"x-client-trace-id",         // Custom trace ID header for debugging
		"x-temporal-correlation-id", // Cross-cluster correlation ID for linking logs
		"x-request-id",              // Request ID for tracing individual requests
		"x-b3-traceid",              // B3 tracing header for distributed tracing
		"x-b3-sampled",              // B3 sampling decision for tracing
	}
)

// GetValues returns header values for passed header names.
// It always returns slice of the same size as number of passed header names.
func GetValues(ctx context.Context, headerNames ...string) []string {
	headerValues := make([]string, len(headerNames))

	for i, headerName := range headerNames {
		if values := metadata.ValueFromIncomingContext(ctx, headerName); len(values) > 0 {
			headerValues[i] = values[0]
		}
	}

	return headerValues
}

// Propagate propagates version headers from incoming context to outgoing context.
// It copies all headers to outgoing context only if they are exist in incoming context
// and doesn't exist in outgoing context already.
func Propagate(ctx context.Context) context.Context {
	headersToAppend := make([]string, 0, len(propagateHeaders)*2)
	mdOutgoing, mdOutgoingExist := metadata.FromOutgoingContext(ctx)
	propagatedHeaders := make(map[string]string)

	for _, headerName := range propagateHeaders {
		if incomingValue := metadata.ValueFromIncomingContext(ctx, headerName); len(incomingValue) > 0 && len(mdOutgoing.Get(headerName)) == 0 {
			headersToAppend = append(headersToAppend, headerName, incomingValue[0])
			propagatedHeaders[headerName] = incomingValue[0]
		}
	}

	// Log header propagation for tracing
	if len(propagatedHeaders) > 0 {
		// Extract request ID for correlation
		requestID := propagatedHeaders["x-request-id"]
		if requestID == "" {
			requestID = "unknown"
		}

		// Log the propagated headers
		headerList := make([]string, 0, len(propagatedHeaders))
		for k, v := range propagatedHeaders {
			headerList = append(headerList, fmt.Sprintf("%s=%s", k, v))
		}

		// This will be logged by the service's logger when available
		// For now, we'll use a simple approach that can be enhanced with proper logger injection
		_ = requestID  // Placeholder for future logger integration
		_ = headerList // Placeholder for future logger integration
	}
	if headersToAppend != nil {
		if mdOutgoingExist {
			ctx = metadata.AppendToOutgoingContext(ctx, headersToAppend...)
		} else {
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(headersToAppend...))
		}
	}
	return ctx
}

// HeaderGetter is an interface for getting a single header value from a case insensitive key.
type HeaderGetter interface {
	Get(string) string
}

// Wrapper for gRPC metadata that exposes a helper to extract a single metadata value.
type GRPCHeaderGetter struct {
	ctx context.Context
}

func NewGRPCHeaderGetter(ctx context.Context) GRPCHeaderGetter {
	return GRPCHeaderGetter{ctx: ctx}
}

// Get a single value from the underlying gRPC metadata.
// Returns an empty string if the metadata key is unset.
func (h GRPCHeaderGetter) Get(key string) string {
	if values := metadata.ValueFromIncomingContext(h.ctx, key); len(values) > 0 {
		return values[0]
	}
	return ""
}
