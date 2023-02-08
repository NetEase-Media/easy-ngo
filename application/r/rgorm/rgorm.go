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

package rgorm

import (
	"context"

	conf "github.com/NetEase-Media/easy-ngo/config"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngo/application/r/rmetrics"
	"github.com/NetEase-Media/easy-ngo/clients/xgorm"
	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
)

const (
	key_prefix = "ngo.client.gorm"
)

var dbClients map[string]*xgorm.Client

func init() {
	hooks.Register(hooks.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	opts := make([]xgorm.Option, 0)
	conf.Get(key_prefix, &opts)
	if len(opts) == 0 {
		panic("no gorm config!")
	}
	metrics := rmetrics.GetProvider()
	for _, opt := range opts {
		cli, err := xgorm.New(&opt, &xfmt.XFmt{}, metrics, nil)
		if err != nil {
			panic("init db failed.")
		}
		set(key_prefix+"."+opt.Name, cli)
	}
	return nil
}

func set(key string, client *xgorm.Client) {
	if dbClients == nil {
		dbClients = make(map[string]*xgorm.Client, 1)
	}
	dbClients[key] = client
}

func GetDBClient(key string) (cli *xgorm.Client) {
	var ok bool
	cli, ok = dbClients[key]
	if !ok {
		return nil
	}
	return cli
}
