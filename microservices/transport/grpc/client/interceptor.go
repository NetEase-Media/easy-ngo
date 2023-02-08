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
	"context"
	"reflect"

	xgrpc "github.com/NetEase-Media/easy-ngo/microservices/transport/grpc"

	"github.com/NetEase-Media/easy-ngoservices/internal/xrpc"
	"github.com/NetEase-Media/easy-ngoservices/middleware"
	"github.com/NetEase-Media/easy-ngoservices/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// unaryClientInterceptor .
func (c *Client) unaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		if c.opts.timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, c.opts.timeout)
			defer cancel()
		}

		op, _ := xrpc.ParseToOperation(method)

		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			if md, ok := xrpc.FromOutgoingContext(ctx); ok {
				ctx = metadata.NewOutgoingContext(ctx, metadata.MD(md))
			}
			err := invoker(ctx, method, req, reply, cc, opts...)
			if in, ok := metadata.FromIncomingContext(ctx); ok {
				ctx = xrpc.NewIncomingContext(ctx, xrpc.MD(in))
			}
			return reply, xgrpc.FromRPCError(err)
		}

		if v, ok := c.mws[""]; ok {
			h = middleware.Chain(v...)(h)
		}

		if v, ok := c.mws["/"]; ok {
			h = middleware.Chain(v...)(h)
		}

		if v, ok := c.mws["/"+op.Pkg]; ok {
			h = middleware.Chain(v...)(h)
		}

		if v, ok := c.mws["/"+op.FullService()]; ok {
			h = middleware.Chain(v...)(h)
		}

		if v, ok := c.mws[method]; ok {
			h = middleware.Chain(v...)(h)
		}

		ctx = transport.NewClientContext(ctx,
			transport.ClientInfo{
				Type:   "GRPC",
				Target: cc.Target(),
				Op:     &op,
			},
		)

		r, err := h(ctx, req)
		if err == nil {
			v := reflect.ValueOf(reply).Elem()
			v.Set(reflect.ValueOf(r).Elem())
		}
		return err
	}
}
