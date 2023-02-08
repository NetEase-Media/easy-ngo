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

package server

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/NetEase-Media/easy-ngo/microservices/middleware"
	"github.com/NetEase-Media/easy-ngo/microservices/sd"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"google.golang.org/grpc"
)

// Option is a function that configures the server.
type Option func(o *options)

// WithName returns an Option that sets the name of the server.
func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

// WithNetwork returns an Option that sets the network of the server.
func WithNetwork(network string) Option {
	return func(o *options) {
		o.network = network
	}
}

// WithAddr returns an Option that sets the address of the server.
func WithAddr(addr string) Option {
	return func(o *options) {
		o.addr = addr
	}
}

// WithListener returns an Option that sets the listener of the server.
func WithListener(listener net.Listener) Option {
	return func(o *options) {
		o.listener = listener
	}
}

// WithRegistrar returns an Option that sets the register of the server.
func WithRegistrar(registrar sd.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}

// WithTimeout returns an Option that sets the timeout of the server.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithMiddlewares returns an Option that sets the middlewares of the server.
func WithMiddlewares(mws ...middleware.Middleware) Option {
	return func(o *options) {
		o.mws = mws
	}
}

// WithGRPCOptions returns an Option that sets the grpc options of the server.
func WithGRPCOptions(opts ...grpc.ServerOption) Option {
	return func(o *options) {
		o.gopts = opts
	}
}

// WithLogger returns an Option that sets the logger of the server.
func WithLogger(logger xlog.Logger) Option {
	return func(o *options) {
		o.log = logger
	}
}

// WithTLSConfig returns an Option that sets the tls config of the server.
func WithTLSConfig(tls *tls.Config) Option {
	return func(o *options) {
		o.tls = tls
	}
}

// WithMetadata returns an Option that sets the metadata of the server.
func WithMetadata(metadata map[string]string) Option {
	return func(o *options) {
		o.metadata = metadata
	}
}

// defaultOptions returns the default options.
func defaultOptions() *options {
	return &options{
		name:    "ngo",
		network: "tcp",
		addr:    ":0",
		timeout: time.Second * 30,
		log:     xlog.NewNopLogger(),
	}
}

// options is the options for server.
type options struct {
	name      string
	network   string
	addr      string
	listener  net.Listener
	timeout   time.Duration
	registrar sd.Registrar
	mws       []middleware.Middleware
	tls       *tls.Config
	gopts     []grpc.ServerOption
	log       xlog.Logger
	metadata  map[string]string
}
