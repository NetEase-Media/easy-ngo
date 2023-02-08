// Copyright 2022 NetEase Media Technology（Beijing）Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"context"
	"time"

	"github.com/NetEase-Media/easy-ngo/microservices/errors"
	"github.com/NetEase-Media/easy-ngoservices/middleware"
	"github.com/NetEase-Media/easy-ngoservices/transport"
	"github.com/NetEase-Media/easy-ngovability/metrics"
	"github.com/NetEase-Media/easy-ngo
)

// Option is a function that configures the middleware.
type Option func(*options)

// defaultOptions returns the default options.
func defaultOptions() *options {
	return &options{
		log: xlog.NewNopLogger(),
	}
}

// options is the options for metrics middleware.
type options struct {
	log xlog.Logger
}

// WithLogger returns a new Option that sets the logger.
func WithLogger(logger xlog.Logger) Option {
	return func(o *options) {
		o.log = logger
	}
}

// Server returns a server metrics middleware.
func Server(provider metrics.Provider, opts ...Option) middleware.Middleware {
	opt := defaultOptions()
	for _, o := range opts {
		o(opt)
	}
	counter := provider.NewCounter("rpc_server_requests_total", "type", "operation", "code", "reason")
	histogram := provider.NewHistogram("rpc_server_responses_seconds", []float64{}, "type", "operation")

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			start := time.Now()
			defer func() {
				info, _ := transport.FromServerContext(ctx)
				s := errors.FromError(err)
				counter.With(
					"type", info.Type,
					"operation", info.Op.String(),
					"code", s.Code.String(),
					"reason", s.Reason,
				).Inc()
				histogram.With(
					"type", info.Type,
					"operation", info.Op.String(),
				).Observe(time.Since(start).Seconds())
			}()

			resp, err = handler(ctx, req)
			return
		}
	}
}

// Client returns a client metrics middleware.
func Client(provider metrics.Provider, opts ...Option) middleware.Middleware {
	opt := defaultOptions()
	for _, o := range opts {
		o(opt)
	}

	counter := provider.NewCounter("rpc_client_requests_total", "type", "target", "operation", "code",
		"reason")
	histogram := provider.NewHistogram("rpc_client_responses_seconds", []float64{}, "type", "target",
		"operation")

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			start := time.Now()
			defer func() {
				info, _ := transport.FromClientContext(ctx)
				s := errors.FromError(err)
				counter.With(
					"type", info.Type,
					"target", info.Target,
					"operation", info.Op.String(),
					"code", s.Code.String(),
					"reason", s.Reason,
				).Inc()
				histogram.With(
					"type", info.Type,
					"target", info.Target,
					"operation", info.Op.String(),
				).Observe(time.Since(start).Seconds())
			}()

			resp, err = handler(ctx, req)
			return
		}
	}
}
