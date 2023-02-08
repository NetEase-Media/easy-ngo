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
	"crypto/tls"
	"fmt"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngo/application/r/rlog"
	"github.com/NetEase-Media/easy-ngo/application/r/rmetrics"
	"github.com/NetEase-Media/easy-ngo/application/r/rms"
	"github.com/NetEase-Media/easy-ngo/application/r/rms/api"
	"github.com/NetEase-Media/easy-ngo/application/r/rms/sd"
	"github.com/NetEase-Media/easy-ngo/microservices/middleware"
	"github.com/NetEase-Media/easy-ngo/microservices/middleware/metrics"
	"github.com/NetEase-Media/easy-ngo/microservices/middleware/tracing"
	"github.com/NetEase-Media/easy-ngo/microservices/transport/grpc/server"
	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
)

var servers map[string]*server.Server

func init() {
	hooks.Register(hooks.Initialize, initialize)
}

func initialize(ctx context.Context) error {
	gconfig, err := rms.GetConfig()
	if err != nil {
		return fmt.Errorf("[microservices] load config failed: %v", err)
	}
	config := gconfig.GetGrpc().GetServers()
	servers = make(map[string]*server.Server, len(config))
	for i := range config {
		if _, ok := servers[config[i].Name]; ok {
			return fmt.Errorf("[microservices] duplicate server key: %s", config[i].Name)
		}

		s, err := newServer(ctx, config[i])
		if err != nil {
			return fmt.Errorf("[microservices] new server failed: %v", err)
		}
		servers[config[i].Name] = s
		hooks.Register(hooks.Start, func(ctx context.Context) error {
			return s.Start()
		})
		hooks.Register(hooks.Stop, func(ctx context.Context) error {
			s.GracefulStop()
			return nil
		})
		hooks.Register(hooks.Online, func(ctx context.Context) error {
			return s.Online(ctx)
		})
		hooks.Register(hooks.Offline, func(ctx context.Context) error {
			return s.Offline(ctx)
		})
		hooks.Register(hooks.HealthCheck, func(ctx context.Context) error {
			ok := s.Healthz(ctx)
			if !ok {
				return fmt.Errorf("[microservices] server %s not health", config[i].Name)
			}
			return nil
		})
	}
	return nil
}

func newServer(ctx context.Context, config *api.GRPCServer) (*server.Server, error) {
	opts := make([]server.Option, 0, 10)

	if config.Name != "" {
		opts = append(opts, server.WithName(config.Name))
	}
	if config.Network != "" {
		opts = append(opts, server.WithNetwork(config.Network))
	}
	if config.Addr != "" {
		opts = append(opts, server.WithAddr(config.Addr))
	}
	if config.Timeout != nil && config.Timeout.AsDuration() > 0 {
		opts = append(opts, server.WithTimeout(config.Timeout.AsDuration()))
	}
	if config.RegistrarRef != "" {
		registrar := sd.GetRegistrar(config.RegistrarRef)
		if registrar == nil {
			return nil, fmt.Errorf("[microservices] registrar not found: %s", config.RegistrarRef)
		}
		opts = append(opts, server.WithRegistrar(registrar))
	}
	var log xlog.Logger
	log, _ = xfmt.Default()
	if config.LoggerRef != "" {
		log = rlog.GetLogger(config.LoggerRef)
		if log == nil {
			return nil, fmt.Errorf("[microservices] server logger not found: %s", config.LoggerRef)
		}
	}

	opts = append(opts, server.WithLogger(log))

	if config.EnableMetrics || config.EnableTracing || config.EnableLogging {
		mws := make([]middleware.Middleware, 0, 3)
		if config.EnableMetrics {
			mws = append(mws, metrics.Server(rmetrics.GetProvider(), metrics.WithLogger(log)))
		}
		if config.EnableTracing {
			mws = append(mws, tracing.Server(tracer.GetTracerProvider(), tracing.WithLogger(log)))
		}
		if config.EnableLogging {
			return nil, fmt.Errorf("[microservices] server logging not supported")
		}

		opts = append(opts, server.WithMiddlewares(mws...))
	}

	if config.Tls != nil {
		cert, err := tls.LoadX509KeyPair(config.Tls.CertFile, config.Tls.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("[microservices] server tls error: %s", err)
		}
		tlsConf := &tls.Config{Certificates: []tls.Certificate{cert}}
		opts = append(opts, server.WithTLSConfig(tlsConf))
	}

	if config.GrpcOpts != nil {
		opts = append(opts, server.WithGRPCOptions(getGRPCOptions(config.GrpcOpts)...))
	}

	if len(config.Metadata) > 0 {
		opts = append(opts, server.WithMetadata(config.Metadata))
	}
	return server.New(opts...)
}

func Get(name string) *server.Server {
	return servers[name]
}
