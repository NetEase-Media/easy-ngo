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

package tracing

import (
	"context"
	"strings"

	"github.com/NetEase-Media/easy-ngo/microservices/errors"
	"github.com/NetEase-Media/easy-ngo/microservices/internal/xnet"
	"github.com/NetEase-Media/easy-ngo/microservices/internal/xrpc"
	"github.com/NetEase-Media/easy-ngo/microservices/middleware"
	"github.com/NetEase-Media/easy-ngo/microservices/transport"
	"github.com/NetEase-Media/easy-ngo/observability/tracing"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// Option is a function that configures the middleware.
type Option func(*options)

// defaultOptions returns the default options.
func defaultOptions() *options {
	return &options{
		log: xlog.NewNopLogger(),
	}
}

// options is the options for tracing middleware.
type options struct {
	log xlog.Logger
}

// WithLogger returns a new Option that sets the logger.
func WithLogger(logger xlog.Logger) Option {
	return func(o *options) {
		o.log = logger
	}
}

// Server returns a server tracing middleware.
func Server(provider tracing.Provider, opts ...Option) middleware.Middleware {
	opt := defaultOptions()
	for _, o := range opts {
		o(opt)
	}

	propagator := tracing.GetTextMapPropagator()
	tr := provider.Tracer("RPCServer")

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			info, _ := transport.FromServerContext(ctx)
			md, _ := xrpc.FromIncomingContext(ctx)
			if md == nil {
				md = xrpc.NewMD(nil)
			} else {
				md = md.Copy()
			}

			ctx = propagator.Extract(ctx, MDCarrier(md))
			var span trace.Span
			ctx, span = tr.Start(ctx, info.Op.String(), trace.WithSpanKind(trace.SpanKindServer))

			service, method := info.Op.FullService(), info.Op.Method
			host, port := xnet.ParseAddr(info.Peer.Addr.String())

			span.SetAttributes(
				semconv.RPCSystemKey.String(info.Type),
				semconv.RPCServiceKey.String(service),
				semconv.RPCMethodKey.String(method),
				semconv.NetPeerIPKey.String(host),
				semconv.NetPeerPortKey.String(port),
			)

			defer func() {
				s := errors.FromError(err)
				span.SetAttributes(attribute.Key("rpc.status_code").String(s.Code.String()))
				span.SetAttributes(attribute.Key("rpc.reason").String(s.Reason))
				span.SetAttributes(attribute.Key("rpc.message").String(s.Message))

				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				} else {
					span.SetStatus(codes.Ok, codes.Ok.String())
				}
				span.End()
			}()

			resp, err = handler(ctx, req)
			return
		}
	}
}

// Client returns a client tracing middleware.
func Client(provider tracing.Provider, opts ...Option) middleware.Middleware {
	opt := defaultOptions()
	for _, o := range opts {
		o(opt)
	}

	propagator := tracing.GetTextMapPropagator()
	tr := provider.Tracer("RPCClient")

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			info, _ := transport.FromClientContext(ctx)
			md, _ := xrpc.FromOutgoingContext(ctx)
			if md == nil {
				md = xrpc.NewMD(nil)
				ctx = xrpc.NewOutgoingContext(ctx, md)
			}

			var span trace.Span
			ctx, span = tr.Start(ctx, info.Op.String(), trace.WithSpanKind(trace.SpanKindClient))
			propagator.Inject(ctx, MDCarrier(md))

			service, method := info.Op.FullService(), info.Op.Method

			var host, port string
			if addr, err := xnet.ParseURL(info.Target); err == nil {
				host, port = xnet.ParseAddr(addr)
			}

			span.SetAttributes(
				semconv.RPCSystemKey.String(info.Type),
				semconv.RPCServiceKey.String(service),
				semconv.RPCMethodKey.String(method),
				semconv.NetPeerIPKey.String(host),
				semconv.NetPeerPortKey.String(port),
			)

			defer func() {
				s := errors.FromError(err)
				span.SetAttributes(attribute.Key("rpc.status_code").String(s.Code.String()))
				span.SetAttributes(attribute.Key("rpc.reason").String(s.Reason))
				span.SetAttributes(attribute.Key("rpc.message").String(s.Message))

				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				} else {
					span.SetStatus(codes.Ok, codes.Ok.String())
				}
				span.End()
			}()

			resp, err = handler(ctx, req)
			return
		}
	}
}

// MDCarrier is a carrier for the TextMapPropagator.
type MDCarrier xrpc.MD

var _ propagation.TextMapCarrier = (*MDCarrier)(nil)

// Get returns the value associated with the passed key.
func (c MDCarrier) Get(key string) string {
	vals := xrpc.MD(c).Get(key)
	if len(vals) > 0 {
		return strings.Join(vals, ",")
	}
	return ""
}

// Set stores the key-value pair.
func (c MDCarrier) Set(key, value string) {
	xrpc.MD(c).Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (c MDCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range xrpc.MD(c) {
		keys = append(keys, k)
	}
	return keys
}
