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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/NetEase-Media/easy-ngo/application/r/rlog"
	"github.com/NetEase-Media/easy-ngo

	"github.com/NetEase-Media/easy-ngocation/hooks"
	"github.com/NetEase-Media/easy-ngocation/r/rmetrics"
	"github.com/NetEase-Media/easy-ngocation/r/rms"
	"github.com/NetEase-Media/easy-ngocation/r/rms/api"
	"github.com/NetEase-Media/easy-ngocation/r/rms/sd"
	"github.com/NetEase-Media/easy-ngoservices/middleware"
	"github.com/NetEase-Media/easy-ngoservices/middleware/metrics"
	"github.com/NetEase-Media/easy-ngoservices/middleware/tracing"
	"github.com/NetEase-Media/easy-ngoservices/transport/grpc/client"
	tracer "github.com/NetEase-Media/easy-ngovability/tracing"
	"github.com/NetEase-Media/easy-ngoxfmt"
)

var clients map[string]*client.Client

func init() {
	hooks.Register(hooks.Initialize, initialize)
}

func initialize(ctx context.Context) error {
	gconfig, err := rms.GetConfig()
	if err != nil {
		return fmt.Errorf("[microservices] load config failed: %v", err)
	}
	config := gconfig.GetGrpc().GetClients()
	clients = make(map[string]*client.Client, len(config))
	for i := range config {
		if _, ok := clients[config[i].Name]; ok {
			return fmt.Errorf("[microservices] duplicate client key: %s", config[i].Name)
		}

		c, err := newClient(ctx, config[i])
		if err != nil {
			return fmt.Errorf("[microservices] new client failed: %v", err)
		}
		clients[config[i].Name] = c
		hooks.Register(hooks.Stop, func(ctx context.Context) error {
			c.Close()
			return nil
		})
	}
	return nil
}

func newClient(ctx context.Context, config *api.GRPCClient) (*client.Client, error) {
	target := config.Target
	if target == "" {
		return nil, fmt.Errorf("[microservices] client target must be set")
	}

	opts := make([]client.Option, 0, 10)

	if config.Timeout != nil && config.Timeout.AsDuration() > 0 {
		opts = append(opts, client.WithTimeout(config.Timeout.AsDuration()))
	}
	if config.DiscoveryRef != "" {
		opts = append(opts, client.WithDiscovery(sd.GetDiscovery(config.DiscoveryRef)))
	}
	if config.BalancerName != "" {
		opts = append(opts, client.WithBalancerName(config.BalancerName))
	}
	opts = append(opts, client.WithEnabledHealthCheck(config.EnabledHealthCheck))

	var log xlog.Logger
	log, _ = xfmt.Default()
	if config.LoggerRef != "" {
		log = rlog.GetLogger(config.LoggerRef)
		if log == nil {
			return nil, fmt.Errorf("[microservices] client logger not found: %s", config.LoggerRef)
		}
	}

	opts = append(opts, client.WithLogger(log))

	if config.EnableMetrics || config.EnableTracing || config.EnableLogging {
		mws := make([]middleware.Middleware, 0, 3)
		if config.EnableMetrics {
			mws = append(mws, metrics.Client(rmetrics.GetProvider(), metrics.WithLogger(log)))
		}
		if config.EnableTracing {
			mws = append(mws, tracing.Client(tracer.GetTracerProvider(), tracing.WithLogger(log)))
		}
		if config.EnableLogging {
			return nil, fmt.Errorf("[microservices] client logging not supported")
		}

		opts = append(opts, client.WithMiddlewares(mws...))
	}

	if config.Tls != nil {
		b, err := os.ReadFile(config.Tls.CertFile)
		if err != nil {
			return nil, fmt.Errorf("[microservices] client tls error: %s", err)
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			return nil, fmt.Errorf("[microservices] client tls error: %s", err)
		}
		tlsConf := &tls.Config{ServerName: config.Tls.ServerName, RootCAs: cp}
		opts = append(opts, client.WithTLSConfig(tlsConf))
	}

	if config.GrpcOpts != nil {
		opts = append(opts, client.WithGRPCOptions(getGRPCOptions(config.GrpcOpts)...))
	}

	return client.New(ctx, target, opts...)
}

func Get(name string) *client.Client {
	return clients[name]
}
