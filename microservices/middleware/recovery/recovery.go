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

package recovery

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/microservices/errors"
	"github.com/NetEase-Media/easy-ngoservices/middleware"
	"github.com/NetEase-Media/easy-ngo
)

// Handler is a function that handles a panic.
type Handler func(ctx context.Context, req, p interface{}) error

// Option is a function that configures the middleware.
type Option func(*options)

// defaultOptions returns the default options.
func defaultOptions() *options {
	return &options{
		handler: func(ctx context.Context, req, p interface{}) error {
			return errors.Newf(errors.Internal, "INTERNAL_ERROR", "%v", p)
		},
		log: xlog.NewNopLogger(),
	}
}

// options is the options for recovery middleware.
type options struct {
	handler Handler
	log     xlog.Logger
}

// WithHandler returns a new Option that sets the handler.
func WithHandler(h Handler) Option {
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

// Recovery returns a recovery middleware.
func Recovery(opts ...Option) middleware.Middleware {
	opt := defaultOptions()
	for _, o := range opts {
		o(opt)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					opt.log.Errorf("internal error r=%v", r)
					err = opt.handler(ctx, req, r)
				}
			}()
			resp, err = handler(ctx, req)
			return
		}
	}
}
