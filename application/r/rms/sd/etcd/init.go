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

package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/NetEase-Media/easy-ngo/application/r/rlog"
	"github.com/NetEase-Media/easy-ngo

	"github.com/NetEase-Media/easy-ngocation/hooks"
	"github.com/NetEase-Media/easy-ngocation/r/rms"
	"github.com/NetEase-Media/easy-ngocation/r/rms/api"
	"github.com/NetEase-Media/easy-ngocation/r/rms/sd/internal"
	"github.com/NetEase-Media/easy-ngoservices/contrib/sd/etcd"
	"github.com/NetEase-Media/easy-ngoservices/sd"
	"github.com/NetEase-Media/easy-ngoxfmt"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func init() {
	hooks.Register(hooks.Initialize, initialize)
}

func initialize(ctx context.Context) error {
	gconfig, err := rms.GetConfig()
	if err != nil {
		return fmt.Errorf("[microservices] load config failed: %v", err)
	}
	config := gconfig.GetSd().GetEtcds()
	if internal.Discoveries == nil {
		internal.Discoveries = make(map[string]sd.Discovery, len(config))
	}
	if internal.Registrars == nil {
		internal.Registrars = make(map[string]sd.Registrar, len(config))
	}

	for i := range config {
		if _, ok := internal.Discoveries[config[i].Name]; ok {
			return fmt.Errorf("[microservices] duplicate discovery key: %s", config[i].Name)
		}
		if _, ok := internal.Registrars[config[i].Name]; ok {
			return fmt.Errorf("[microservices] duplicate registrar key: %s", config[i].Name)
		}

		c, err := newSD(ctx, config[i])
		if err != nil {
			return fmt.Errorf("[microservices] new sd failed: %v", err)
		}
		internal.Discoveries[config[i].Name] = c
		internal.Registrars[config[i].Name] = c
	}
	return nil
}

func newSD(ctx context.Context, config *api.Etcd) (*etcd.ServiceDiscovery, error) {
	if len(config.Endpoints) == 0 {
		return nil, fmt.Errorf("[microservices] sd etcd endpoints must be set")
	}

	conf := clientv3.Config{
		Endpoints:            config.Endpoints,
		DialTimeout:          10 * time.Second,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		DialOptions:          []grpc.DialOption{grpc.WithBlock()},
	}
	if config.ConnectTimeout != nil && config.ConnectTimeout.AsDuration() > 0 {
		conf.DialTimeout = config.ConnectTimeout.AsDuration()
	}
	if config.Tls != nil {
		b, err := os.ReadFile(config.Tls.CertFile)
		if err != nil {
			return nil, fmt.Errorf("[microservices] sd etcd tls error: %s", err)
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			return nil, fmt.Errorf("[microservices] sd etcd tls error: %s", err)
		}
		tlsConf := &tls.Config{ServerName: config.Tls.ServerName, RootCAs: cp}
		conf.TLS = tlsConf
	}
	if config.Auth != nil {
		conf.Username = config.Auth.Username
		conf.Password = config.Auth.Password
	}
	etcdClient, err := clientv3.New(conf)
	if err != nil {
		return nil, fmt.Errorf("[microservices] new etcd client failed: %v", err)
	}
	opts := make([]etcd.Option, 0, 3)
	if config.Namespace != "" {
		opts = append(opts, etcd.WithNamespace(config.Namespace))
	}
	if config.Ttl != nil && config.Ttl.AsDuration() > 0 {
		opts = append(opts, etcd.WithTTL(config.Ttl.AsDuration()))
	}

	var log xlog.Logger
	log, _ = xfmt.Default()
	if config.LoggerRef != "" {
		log = rlog.GetLogger(config.LoggerRef)
		if log == nil {
			return nil, fmt.Errorf("[microservices] sd logger not found: %s", config.LoggerRef)
		}
	}

	opts = append(opts, etcd.WithLogger(log))

	return etcd.New(ctx, etcdClient, opts...), nil
}
