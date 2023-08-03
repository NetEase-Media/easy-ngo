// Copyright 2023 NetEase Media Technology（Beijing）Co., Ltd.
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

package rgoroutinepool

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngo/clients/xgoroutinepool"
	conf "github.com/NetEase-Media/easy-ngo/config"
)

const (
	key_prefix = "ngo.client.routinepool"
)

var routinePool map[string]xgoroutinepool.Pool

func init() {
	hooks.Register(hooks.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	opts := make([]xgoroutinepool.Option, 0)
	err := conf.Get(key_prefix, &opts)
	if err != nil {
		panic(err)
	}
	for _, opt := range opts {
		cli := xgoroutinepool.New(&opt)
		if err != nil {
			panic("init go routine pool failed.")
		}
		k := key_prefix + "." + opt.Name
		set(k, cli)
	}
	return nil
}

func set(key string, client xgoroutinepool.Pool) {
	if routinePool == nil {
		routinePool = make(map[string]xgoroutinepool.Pool, 1)
	}
	routinePool[key] = client
}

func GetHttpClient(key string) (cli xgoroutinepool.Pool) {
	var ok bool
	cli, ok = routinePool[key]
	if !ok {
		return nil
	}
	return cli
}
