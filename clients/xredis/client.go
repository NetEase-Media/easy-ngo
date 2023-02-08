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

package xredis

import (
	"github.com/NetEase-Media/easy-ngo/observability/metrics"
	tracer "github.com/NetEase-Media/easy-ngovability/tracing"
	"github.com/NetEase-Media/easy-ngo
	"github.com/go-redis/redis/v8"
)

func newClientOptions(opt *Option) *redis.Options {
	return &redis.Options{
		Addr:               opt.Addr[0],
		Username:           opt.Username,
		Password:           opt.Password,
		DB:                 opt.DB,
		MaxRetries:         opt.MaxRetries,
		MinRetryBackoff:    opt.MinRetryBackoff,
		MaxRetryBackoff:    opt.MaxRetryBackoff,
		DialTimeout:        opt.DialTimeout,
		ReadTimeout:        opt.ReadTimeout,
		WriteTimeout:       opt.WriteTimeout,
		PoolSize:           opt.PoolSize,
		MinIdleConns:       opt.MinIdleConns,
		MaxConnAge:         opt.MaxConnAge,
		PoolTimeout:        opt.PoolTimeout,
		IdleTimeout:        opt.IdleTimeout,
		IdleCheckFrequency: opt.IdleCheckFrequency,
		TLSConfig:          opt.TLSConfig,
	}
}

func NewClient(opt *Option, logger xlog.Logger, metrics metrics.Provider, tracer tracer.Provider) *RedisContainer {
	baseClient := redis.NewClient(newClientOptions(opt))
	c := &RedisContainer{
		Redis:     baseClient,
		Opt:       *opt,
		redisType: RedisTypeClient,
		logger:    logger,
		metrics:   metrics,
		tracer:    tracer,
	}

	if metrics != nil {
		metricsHook := newMetricHook(c, logger, metrics)
		baseClient.AddHook(metricsHook)
	}
	// baseClient.AddHook(newMetricHook(c))
	if opt.EnableTracer {
		baseClient.AddHook(newTracingHook(c))
	}
	return c
}
