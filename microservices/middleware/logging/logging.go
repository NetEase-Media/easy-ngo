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

package logging

import (
	"context"
	"time"

	"github.com/NetEase-Media/easy-ngo/microservices/errors"
	"github.com/NetEase-Media/easy-ngo/microservices/internal/xnet"
	"github.com/NetEase-Media/easy-ngo/microservices/middleware"
	"github.com/NetEase-Media/easy-ngo/microservices/transport"
	"github.com/NetEase-Media/easy-ngo/observability/logging"
	"github.com/NetEase-Media/easy-ngo/xlog"
)

// Option is a function that configures the middleware.
type Option func(*options)

// defaultOptions returns the default options.
func defaultOptions() *options {
	return &options{
		log: xlog.NewNopLogger(),
	}
}

// options is the options for logging middleware.
type options struct {
	log xlog.Logger
}

// WithLogger returns a new Option that sets the logger.
func WithLogger(logger xlog.Logger) Option {
	return func(o *options) {
		o.log = logger
	}
}

// Server returns a server logging middleware.
func Server(provider logging.Logger, opts ...Option) middleware.Middleware {
	opt := defaultOptions()
	for _, o := range opts {
		o(opt)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			start := time.Now()
			defer func() {
				info, _ := transport.FromServerContext(ctx)
				service, method := info.Op.FullService(), info.Op.Method
				host, port := xnet.ParseAddr(info.Peer.Addr.String())

				s := errors.FromError(err)
				location, _ := time.LoadLocation("Local")
				provider.Log(
					"timestamp", start.In(location),
					"type", info.Type,
					"service", service,
					"method", method,
					"host", host,
					"port", port,
					"code", s.Code.String(),
					"reason", s.Reason,
					"latency", time.Since(start).Seconds(),
				)
			}()

			resp, err = handler(ctx, req)
			return
		}
	}
}

// Client returns a client logging middleware.
func Client(provider logging.Logger, opts ...Option) middleware.Middleware {
	opt := defaultOptions()
	for _, o := range opts {
		o(opt)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			start := time.Now()
			defer func() {
				info, _ := transport.FromClientContext(ctx)
				service, method := info.Op.FullService(), info.Op.Method

				var host, port string
				if addr, err := xnet.ParseURL(info.Target); err == nil {
					host, port = xnet.ParseAddr(addr)
				}

				s := errors.FromError(err)
				location, _ := time.LoadLocation("Local")

				provider.Log(
					"timestamp", start.In(location),
					"type", info.Type,
					"service", service,
					"method", method,
					"host", host,
					"port", port,
					"code", s.Code.String(),
					"reason", s.Reason,
					"latency", time.Since(start).Seconds(),
				)
			}()

			return handler(ctx, req)
		}
	}
}
