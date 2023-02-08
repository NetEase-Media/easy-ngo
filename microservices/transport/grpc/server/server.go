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
	"errors"
	"net"
	"strconv"

	"google.golang.org/grpc/credentials"

	"github.com/NetEase-Media/easy-ngo/microservices/internal/xnet"
	"github.com/NetEase-Media/easy-ngo/microservices/middleware"
	"github.com/NetEase-Media/easy-ngo/microservices/sd"
	"github.com/NetEase-Media/easy-ngo/microservices/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	healthServiceName = "grpc_health_v1"
)

var _ transport.Server = (*Server)(nil)

var (
	UnavailableAddr = errors.New("unavailable addr")
	UnableToGetPort = errors.New("unable to get port")
)

// New creates a new grpc server.
func New(optFns ...Option) (*Server, error) {
	srv := Server{
		opts:   defaultOptions(),
		mws:    make(map[string][]middleware.Middleware),
		health: health.NewServer(),
	}
	for _, o := range optFns {
		o(srv.opts)
	}

	if len(srv.opts.mws) > 0 {
		srv.mws[""] = srv.opts.mws
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(srv.unaryServerInterceptor()),
	}
	if srv.opts.tls != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.opts.tls)))
	}
	grpcOpts = append(grpcOpts, srv.opts.gopts...)
	srv.Server = grpc.NewServer(grpcOpts...)
	srv.health.SetServingStatus(healthServiceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	reflection.Register(srv.Server)
	return &srv, nil
}

// Server is a grpc server.
type Server struct {
	*grpc.Server
	opts   *options
	addr   string
	lis    net.Listener
	mws    map[string][]middleware.Middleware
	health *health.Server
	err    error
}

// Use adds middleware to the server.
// selector support /, /{package} /{package}.{service}, /{package}.{service}/{method}
func (s *Server) Use(selector string, mws ...middleware.Middleware) {
	if _, ok := s.mws[selector]; !ok {
		s.mws[selector] = make([]middleware.Middleware, 0, len(mws))
	}
	s.mws[selector] = append(s.mws[selector], mws...)
}

// Start starts the server.
func (s *Server) Start() (err error) {
	lis := s.opts.listener
	if lis == nil {
		lis, err = net.Listen(s.opts.network, s.opts.addr)
		if err != nil {
			return err
		}
	}
	ip, err := xnet.LocalIP()
	if err != nil {
		return err
	}
	port, ok := xnet.Port(lis)
	if !ok {
		return UnableToGetPort
	}
	s.addr = net.JoinHostPort(ip, strconv.Itoa(port))
	s.lis = lis
	s.opts.log.Infof("grpc server listening on [:%d]", port)
	return s.Server.Serve(s.lis)
}

// Healthz returns the health status of the server.
func (s *Server) Healthz(ctx context.Context) bool {
	resp, err := s.health.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: healthServiceName})
	if err != nil {
		return false
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return false
	}
	conn, err := s.lis.Accept()
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Online
// sets the health status of the server to SERVING.
// registers the server to the service discovery.
func (s *Server) Online(ctx context.Context) error {
	if s.addr == "" {
		return UnavailableAddr
	}
	s.health.Resume()
	if s.opts.registrar != nil {
		serviceInfo := &sd.ServiceInfo{
			Name:     s.opts.name,
			Scheme:   "grpc",
			Addr:     s.addr,
			Metadata: s.opts.metadata,
		}
		err := s.opts.registrar.Register(ctx, serviceInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// Offline
// sets the health status of the server to NOT_SERVING.
// deregisters the server from the service discovery.
func (s *Server) Offline(ctx context.Context) error {
	s.health.Shutdown()
	if s.opts.registrar != nil {
		serviceInfo := &sd.ServiceInfo{
			Name:     s.opts.name,
			Scheme:   "grpc",
			Addr:     s.addr,
			Metadata: s.opts.metadata,
		}
		err := s.opts.registrar.Deregister(ctx, serviceInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop stops the server.
func (s *Server) Stop() {
	s.Server.Stop()
}

// GracefulStop stops the server gracefully.
func (s *Server) GracefulStop() {
	s.Server.GracefulStop()
}
