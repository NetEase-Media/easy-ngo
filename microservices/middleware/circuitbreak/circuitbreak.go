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

package circuitbreak

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/microservices/transport"
	"github.com/NetEase-Media/easy-ngo/xlog"

	"github.com/NetEase-Media/easy-ngo/microservices/errors"
	"github.com/NetEase-Media/easy-ngo/microservices/middleware"
)

// Option is a function that configures the middleware.
type Option func(*options)

// defaultOptions returns the default options.
func defaultOptions() *options {
	return &options{
		handler: func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, errors.Newf(errors.Unavailable, "CIRCUIT_BREAK", "circuit breaker is open")
		},
		log: xlog.NewNopLogger(),
	}
}

// options is the options for rate limit middleware.
type options struct {
	handler middleware.Handler
	log     xlog.Logger
}

// WithHandler returns a new Option that sets the handler.
func WithHandler(h middleware.Handler) Option {
	return func(o *options) {
		o.handler = h
	}
}

// WithLogger returns a new Option that sets the logger.
func WithLogger(logger xlog.Logger) Option {
	return func(o *options) {
		o.log = logger
	}
}

// CircuitBreak returns a circuit break middleware.
func CircuitBreak(cb CircuitBreaker, opts ...Option) middleware.Middleware {
	opt := defaultOptions()
	for _, o := range opts {
		o(opt)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			isBreak := cb.Execute(func() error {
				resp, err = handler(ctx, req)
				return err
			})
			if isBreak {
				info, _ := transport.FromClientContext(ctx)
				opt.log.Warnf("[%s]: circuit break", info.Op.String())
				resp, err = opt.handler(ctx, req)
			}
			return
		}
	}
}

type CircuitBreaker interface {
	Execute(handler func() (err error)) (isBreak bool)
}
