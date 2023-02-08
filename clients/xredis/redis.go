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
	"errors"
	"io"

	"github.com/NetEase-Media/easy-ngo/observability/metrics"
	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/go-redis/redis/v8"
)

const (
	RedisTypeClient          = "client"
	RedisTypeCluster         = "cluster"
	RedisTypeSentinel        = "sentinel"
	RedisTypeShardedSentinel = "sharded_sentinel"
)

type Redis interface {
	redis.Cmdable
	io.Closer
}

func New(opt *Option, logger xlog.Logger, metrics metrics.Provider, tracer tracer.Provider) (*RedisContainer, error) {
	return newWithOption(opt, logger, metrics, tracer)
}

func newWithOption(opt *Option, logger xlog.Logger, metrics metrics.Provider, tracer tracer.Provider) (*RedisContainer, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}

	var c *RedisContainer
	// 判断连接类型
	switch opt.ConnType {
	case RedisTypeClient:
		c = NewClient(opt, logger, metrics, tracer)
	case RedisTypeCluster:
		c = NewClusterClient(opt, logger, metrics, tracer)
	case RedisTypeSentinel:
		if len(opt.MasterNames) == 0 {
			err := errors.New("empty master name")
			return nil, err
		}
		c = NewSentinelClient(opt, logger, metrics, tracer)
	case RedisTypeShardedSentinel:
		if len(opt.MasterNames) == 0 {
			err := errors.New("empty master name")
			return nil, err
		}
		c = NewShardedSentinelClient(opt, logger, metrics, tracer)
	default:
		err := errors.New("redis connection type need ")
		return nil, err
	}

	return c, nil
}

// RedisContainer 用来存储redis客户端及其额外信息
type RedisContainer struct {
	Redis
	Opt       Option
	redisType string
	logger    xlog.Logger
	metrics   metrics.Provider
	tracer    tracer.Provider
}
