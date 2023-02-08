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

package rhttplib

import (
	"context"

	conf "github.com/NetEase-Media/easy-ngo/config"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngo/application/r/rmetrics"
	"github.com/NetEase-Media/easy-ngo/clients/httplib"
	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
)

const (
	key_prefix = "ngo.client.http"
)

var httpClients map[string]*httplib.HttpClient

func init() {
	hooks.Register(hooks.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	opts := make([]httplib.Option, 0)
	err := conf.Get(key_prefix, &opts)
	if err != nil {
		panic(err)
	}
	logger, _ := xfmt.Default()
	metrics := rmetrics.GetProvider()
	for _, opt := range opts {
		cli, err := httplib.New(&opt, logger, metrics, nil)
		if err != nil {
			panic("init http client failed.")
		}
		k := key_prefix + "." + opt.Name
		set(k, cli)
	}
	return nil
}

func set(key string, client *httplib.HttpClient) {
	if httpClients == nil {
		httpClients = make(map[string]*httplib.HttpClient, 1)
	}
	httpClients[key] = client
}

func GetHttpClient(key string) (cli *httplib.HttpClient) {
	var ok bool
	cli, ok = httpClients[key]
	if !ok {
		return nil
	}
	return cli
}
