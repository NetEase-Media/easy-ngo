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

package xmemcache

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/observability/metrics"
	tracer "github.com/NetEase-Media/easy-ngovability/tracing"
	"github.com/NetEase-Media/easy-ngo
	"github.com/bradfitz/gomemcache/memcache"
)

func New(opt *Option, logger xlog.Logger, metrics metrics.Provider, tracer tracer.Provider) (*MemcacheProxy, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}
	c := memcache.New(opt.Addr...)
	c.Timeout = opt.Timeout
	c.MaxIdleConns = opt.MaxIdleConns
	p := &MemcacheProxy{
		base:    c,
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
		hooks:   make([]Hook, 0),
	}
	if opt.EnableTracer {
		p.AddHook(NewTracingHook())
	}
	if p.metrics != nil {
		initMetrics(p.metrics)
	}
	return p, nil
}

// func NewFromKey(key string) (*MemcacheProxy, error) {
// 	opt := defaultOption()
// 	err := conf.Get(key, opt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return New(opt)
// }

type Hook interface {
	Before(context.Context, string, ...string) (context.Context, error)
	After(context.Context, error)
}

// MemcacheProxy memcache 三方包的包装器类
type MemcacheProxy struct {
	Opt     *Option
	base    *memcache.Client
	logger  xlog.Logger
	metrics metrics.Provider
	tracer  tracer.Provider
	hooks   []Hook
}

func (mp *MemcacheProxy) Initialize() error {
	return nil
}

func (mp *MemcacheProxy) Destory() error {
	return nil
}

func (mp *MemcacheProxy) WithMetrics(metrics metrics.Provider) {
	mp.metrics = metrics
}

// func (mp *MemcacheProxy) WithTracer(tracer tracer.Tracer) {
// 	mp.tracer = tracer
// }

func (mp *MemcacheProxy) WithLogger(logger xlog.Logger) {
	mp.logger = logger
}

func (mp *MemcacheProxy) AddHook(h Hook) {
	mp.hooks = append(mp.hooks, h)
}
