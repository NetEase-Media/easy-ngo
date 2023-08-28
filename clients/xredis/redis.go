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

func New(opt *Config) (*RedisContainer, error) {
	return newWithConfig(opt)
}

func newWithConfig(opt *Config) (*RedisContainer, error) {
	if err := checkConfig(opt); err != nil {
		return nil, err
	}

	var c *RedisContainer
	// 判断连接类型
	switch opt.ConnType {
	case RedisTypeClient:
		c = NewClient(opt)
	case RedisTypeCluster:
		c = NewClusterClient(opt)
	case RedisTypeSentinel:
		if len(opt.MasterNames) == 0 {
			err := errors.New("empty master name")
			return nil, err
		}
		c = NewSentinelClient(opt)
	case RedisTypeShardedSentinel:
		if len(opt.MasterNames) == 0 {
			err := errors.New("empty master name")
			return nil, err
		}
		c = NewShardedSentinelClient(opt)
	default:
		err := errors.New("redis connection type need ")
		return nil, err
	}

	return c, nil
}

// RedisContainer 用来存储redis客户端及其额外信息
type RedisContainer struct {
	Redis
	Opt       Config
	redisType string
}
