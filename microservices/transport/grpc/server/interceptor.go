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
	"context"

	"github.com/NetEase-Media/easy-ngo/microservices/internal/xrpc"
	"github.com/NetEase-Media/easy-ngo/microservices/middleware"
	"github.com/NetEase-Media/easy-ngo/microservices/transport"
	xgrpc "github.com/NetEase-Media/easy-ngo/microservices/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// unaryServerInterceptor .
func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		if s.opts.timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, s.opts.timeout)
			defer cancel()
		}

		op, _ := xrpc.ParseToOperation(info.FullMethod)

		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		if v, ok := s.mws[""]; ok {
			h = middleware.Chain(v...)(h)
		}

		if v, ok := s.mws["/"]; ok {
			h = middleware.Chain(v...)(h)
		}

		if v, ok := s.mws["/"+op.Pkg]; ok {
			h = middleware.Chain(v...)(h)
		}

		if v, ok := s.mws["/"+op.FullService()]; ok {
			h = middleware.Chain(v...)(h)
		}

		if v, ok := s.mws[info.FullMethod]; ok {
			h = middleware.Chain(v...)(h)
		}

		p, _ := peer.FromContext(ctx)
		ctx = transport.NewServerContext(ctx,
			transport.ServerInfo{
				Type: "GRPC",
				Op:   &op,
				Peer: &xrpc.Peer{Addr: p.Addr},
			},
		)

		if in, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = xrpc.NewIncomingContext(ctx, xrpc.MD(in))
		}

		resp, err := h(ctx, req)

		if out, ok := xrpc.FromOutgoingContext(ctx); ok {
			_ = grpc.SetHeader(ctx, metadata.MD(out))
		}

		return resp, xgrpc.ToRPCError(err)
	}
}
