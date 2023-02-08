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
	"fmt"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/NetEase-Media/easy-ngo/microservices/middleware"
	"github.com/NetEase-Media/easy-ngoservices/transport"
	"github.com/NetEase-Media/easy-ngoservices/transport/grpc/client/resolver"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/health"
)

const (
	healthServiceName = "grpc_health_v1"
)

var _ transport.Client = (*Client)(nil)

// New returns a new grpc client.
func New(ctx context.Context, target string, optFns ...Option) (*Client, error) {
	cli := Client{
		opts: defaultOptions(),
		mws:  make(map[string][]middleware.Middleware),
	}
	for _, o := range optFns {
		o(cli.opts)
	}

	if len(cli.opts.mws) > 0 {
		cli.mws[""] = cli.opts.mws
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(cli.unaryClientInterceptor()),
	}
	if cli.opts.tls != nil {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(cli.opts.tls)))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if cli.opts.discovery != nil {
		grpcOpts = append(grpcOpts,
			grpc.WithResolvers(resolver.NewBuilder(cli.opts.discovery)))
	}
	sc := fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, cli.opts.balancerName)
	if cli.opts.enabledHealthCheck {
		sc = fmt.Sprintf(`{"LoadBalancingPolicy": "%s", "HealthCheckConfig": {"ServiceName": "%s"}}`,
			cli.opts.balancerName, healthServiceName)
	}
	grpcOpts = append(grpcOpts, grpc.WithDefaultServiceConfig(sc))

	if len(cli.opts.gopts) > 0 {
		grpcOpts = append(grpcOpts, cli.opts.gopts...)
	}

	conn, err := grpc.DialContext(ctx, target, grpcOpts...)
	cli.ClientConn = conn

	return &cli, err
}

// Client is a grpc client.
type Client struct {
	*grpc.ClientConn
	opts *options
	mws  map[string][]middleware.Middleware
}

// Use adds middlewares to the client.
// selector support /, /{package} /{package}.{service}, /{package}.{service}/{method}
func (c *Client) Use(selector string, mws ...middleware.Middleware) {
	if _, ok := c.mws[selector]; !ok {
		c.mws[selector] = make([]middleware.Middleware, 0, len(mws))
	}
	c.mws[selector] = append(c.mws[selector], mws...)
}

//func (c *Client) Invoke(ctx context.Context, method string, req interface{}, resp interface{}, opts ...transport.CallOption) error {
//	return c.ClientConn.Invoke(ctx, method, req, resp)
//}

// Close closes the client.
func (c *Client) Close() error {
	return c.ClientConn.Close()
}
