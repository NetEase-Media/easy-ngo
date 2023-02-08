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

package ratelimit

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/microservices/errors"
	"github.com/NetEase-Media/easy-ngo/microservices/middleware"
	"github.com/NetEase-Media/easy-ngo/microservices/transport"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"golang.org/x/time/rate"
)

// Handler is a function that handles a limit.
type Handler func(ctx context.Context, req interface{}) error

// Option is a function that configures the middleware.
type Option func(*options)

// defaultOptions returns the default options.
func defaultOptions() *options {
	return &options{
		handler: func(ctx context.Context, req interface{}) error {
			return errors.Newf(errors.ResourceExhausted, "RATE_LIMIT", "too many requests")
		},
		log: xlog.NewNopLogger(),
	}
}

// options is the options for rate limit middleware.
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

// RateLimit returns a rate limit middleware.
func RateLimit(rl RateLimiter, opts ...Option) middleware.Middleware {
	opt := defaultOptions()
	for _, o := range opts {
		o(opt)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			if !rl.Acquire() {
				info, _ := transport.FromServerContext(ctx)
				opt.log.Warnf("[%s]: rate limit exceeded", info.Op.String())
				err = opt.handler(ctx, req)
				return
			}
			resp, err = handler(ctx, req)
			return
		}
	}
}

type RateLimiter interface {
	// Acquire returns true if the rate limit is not exceeded.
	Acquire() bool
}

func NewTokenBucketRateLimiter(limit rate.Limit, b int) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		Limiter: rate.NewLimiter(limit, b),
	}
}

type TokenBucketRateLimiter struct {
	*rate.Limiter
}

func (rl *TokenBucketRateLimiter) Acquire() bool {
	return rl.Limiter.Allow()
}
