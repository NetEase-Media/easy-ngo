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

package rkafka

import (
	"context"

	conf "github.com/NetEase-Media/easy-ngo/config"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngo/application/r/rmetrics"
	"github.com/NetEase-Media/easy-ngo/clients/xkafka"
	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
)

const (
	key = "ngo.client.kafka"
)

var (
	producers = make(map[string]*xkafka.Producer, 1)
	consumers = make(map[string]*xkafka.Consumer, 1)
)

func init() {
	hooks.Register(hooks.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	opts := make([]xkafka.Option, 0)
	err := conf.Get(key, &opts)
	if err != nil {
		panic(err)
	}
	metrics := rmetrics.GetProvider()
	for _, opt := range opts {
		cli, err := xkafka.New(&opt, &xfmt.XFmt{}, metrics, nil)
		if err != nil {
			panic("init kafka failed." + err.Error())
		}
		tk := key + "." + opt.Name
		producers[tk] = cli.Producer
		consumers[tk] = cli.Consumer
	}
	return nil
}

func GetProducer(key string) *xkafka.Producer {
	return producers[key]
}

func GetConsumer(key string) *xkafka.Consumer {
	return consumers[key]
}
