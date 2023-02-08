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

package rmemcache

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngocation/r/rmetrics"
	"github.com/NetEase-Media/easy-ngots/xmemcache"
	conf "github.com/NetEase-Media/easy-ngog"
	"github.com/NetEase-Media/easy-ngoxfmt"
)

const (
	key = "ngo.client.memcache"
)

var memecacheClients map[string]*xmemcache.MemcacheProxy

func init() {
	hooks.Register(hooks.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	opts := make([]xmemcache.Option, 0)
	conf.Get(key, &opts)
	if len(opts) == 0 {
		panic("no gorm config!")
	}
	metrics := rmetrics.GetProvider()
	for _, opt := range opts {
		cli, err := xmemcache.New(&opt, &xfmt.XFmt{}, metrics, nil)
		if err != nil {
			panic("init xmemcache failed.")
		}
		set(key+"."+opt.Name, cli)
	}
	return nil
}

func set(key string, client *xmemcache.MemcacheProxy) {
	if memecacheClients == nil {
		memecacheClients = make(map[string]*xmemcache.MemcacheProxy, 1)
	}
	memecacheClients[key] = client
}

func GetClient(key string) (cli *xmemcache.MemcacheProxy) {
	var ok bool
	cli, ok = memecacheClients[key]
	if !ok {
		return nil
	}
	return cli
}
