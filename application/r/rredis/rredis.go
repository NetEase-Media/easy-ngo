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

package rredis

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngocation/r/rmetrics"
	"github.com/NetEase-Media/easy-ngots/xredis"
	conf "github.com/NetEase-Media/easy-ngog"
	"github.com/NetEase-Media/easy-ngoxfmt"
)

const (
	key = "ngo.client.redis"
)

var redisClients map[string]*xredis.RedisContainer

func init() {
	hooks.Register(hooks.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	opts := make([]xredis.Option, 0)
	conf.Get(key, &opts)
	err := conf.Get(key, &opts)
	if err != nil {
		panic(err)
	}
	metrics := rmetrics.GetProvider()
	for _, opt := range opts {
		cli, err := xredis.New(&opt, &xfmt.XFmt{}, metrics, nil)
		if err != nil {
			panic("init redis failed.")
		}
		tk := key + "." + opt.Name
		if len(redisClients) == 0 {
			redisClients = make(map[string]*xredis.RedisContainer, 1)
		}
		redisClients[tk] = cli
	}
	return nil
}

// func init() {
// 	opt := xredis.NewDefaultOptions()
// 	conf.Get(key, opt)
// 	logger := &xfmt.XFmt{}
// 	client, err := xredis.New(opt, logger, nil, nil)
// 	if err != nil {
// 		panic("init client error")
// 	}
// 	if redisClients == nil {
// 		redisClients = make(map[string]*xredis.RedisContainer, 1)
// 		redisClients[key] = client
// 	}
// }

func Get(key string) *xredis.RedisContainer {
	return redisClients[key]
}
