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

package rzk

import (
	"context"

	conf "github.com/NetEase-Media/easy-ngo/config"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngo/clients/xzk"
)

const (
	key = "ngo.client.zookeeper"
)

var zkClients map[string]*xzk.ZookeeperProxy

func init() {
	hooks.Register(hooks.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	opts := make([]xzk.Option, 0)
	err := conf.Get(key, &opts)
	if err != nil {
		panic(err)
	}
	for _, opt := range opts {
		cli, err := xzk.New(&opt)
		if err != nil {
			panic("init zookeeper failed.")
		}
		k := key + "." + opt.Name
		set(k, cli)
	}
	return nil
}

func set(key string, client *xzk.ZookeeperProxy) {
	if zkClients == nil {
		zkClients = make(map[string]*xzk.ZookeeperProxy, 1)
	}
	zkClients[key] = client
}

func GetZKClient(key string) (cli *xzk.ZookeeperProxy) {
	var ok bool
	cli, ok = zkClients[key]
	if !ok {
		return nil
	}
	return cli
}
