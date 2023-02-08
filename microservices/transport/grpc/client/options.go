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

package client

import (
	"crypto/tls"
	"time"

	"github.com/NetEase-Media/easy-ngo/microservices/middleware"
	"github.com/NetEase-Media/easy-ngo/microservices/sd"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"google.golang.org/grpc"
)

// Option is a function that configures the client.
type Option func(o *options)

// WithTimeout returns an Option that sets the timeout of the client.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithBalancerName returns an Option that sets the balancer name of the client.
func WithBalancerName(name string) Option {
	return func(o *options) {
		o.balancerName = name
	}
}

// WithEnabledHealthCheck returns an Option that sets the health check of the client.
func WithEnabledHealthCheck(enabledHealthCheck bool) Option {
	return func(o *options) {
		o.enabledHealthCheck = enabledHealthCheck
	}
}

// WithDiscovery returns an Option that sets the discoverer of the client.
func WithDiscovery(discovery sd.Discovery) Option {
	return func(o *options) {
		o.discovery = discovery
	}
}

// WithMiddlewares returns an Option that sets the middlewares of the client.
func WithMiddlewares(mws ...middleware.Middleware) Option {
	return func(o *options) {
		o.mws = mws
	}
}

// WithGRPCOptions returns an Option that sets the grpc options of the client.
func WithGRPCOptions(opts ...grpc.DialOption) Option {
	return func(o *options) {
		o.gopts = opts
	}
}

// WithLogger returns an Option that sets the logger of the client.
func WithLogger(logger xlog.Logger) Option {
	return func(o *options) {
		o.log = logger
	}
}

// WithTLSConfig returns an Option that sets the tls config of the client.
func WithTLSConfig(tls *tls.Config) Option {
	return func(o *options) {
		o.tls = tls
	}
}

// defaultOptions returns an default options.
func defaultOptions() *options {
	return &options{
		balancerName: "round_robin",
	}
}

// options is the options of the client.
type options struct {
	timeout            time.Duration
	discovery          sd.Discovery
	balancerName       string
	enabledHealthCheck bool
	mws                []middleware.Middleware
	tls                *tls.Config
	gopts              []grpc.DialOption
	log                xlog.Logger
}
