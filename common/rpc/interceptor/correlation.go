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

package interceptor

import (
	"context"

	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/log/tag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type CorrelationInterceptor struct {
	logger log.Logger
}

func NewCorrelationInterceptor(logger log.Logger) *CorrelationInterceptor {
	return &CorrelationInterceptor{
		logger: logger,
	}
}

func (ci *CorrelationInterceptor) Intercept(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	// Extract correlation ID from incoming metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if correlationIDs := md.Get("x-temporal-correlation-id"); len(correlationIDs) > 0 {
			correlationID := correlationIDs[0]
			ci.logger.Info("Incoming request with correlation ID",
				tag.NewStringTag("method", info.FullMethod),
				tag.NewStringTag("correlation_id", correlationID),
				tag.NewStringTag("request_type", "correlation_tracking"))
		}
		
		// Also log the custom trace ID if present
		if traceIDs := md.Get("x-client-trace-id"); len(traceIDs) > 0 {
			traceID := traceIDs[0]
			ci.logger.Info("Incoming request with trace ID",
				tag.NewStringTag("method", info.FullMethod),
				tag.NewStringTag("trace_id", traceID),
				tag.NewStringTag("request_type", "trace_tracking"))
		}
	}

	// Continue with the request
	return handler(ctx, req)
}
