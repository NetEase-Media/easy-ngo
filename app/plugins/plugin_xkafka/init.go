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

package pluginxkafka

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/NetEase-Media/easy-ngo/clients/xkafka"
	"github.com/NetEase-Media/easy-ngo/config"
)

var (
	producers = make(map[string]*xkafka.Producer, 1)
	consumers = make(map[string]*xkafka.Consumer, 1)
)

func init() {
	app.RegisterPlugin(app.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	configs := make([]xkafka.Config, 0)
	if err := config.UnmarshalKey("kafka", configs); err != nil {
		return err
	}
	if len(configs) == 0 {
		configs = append(configs, *xkafka.DefaultConfig())
	}
	for _, opt := range configs {
		cli, err := xkafka.New(&opt)
		if err != nil {
			panic("init kafka failed." + err.Error())
		}
		producers[opt.Name] = cli.Producer
		consumers[opt.Name] = cli.Consumer
	}
	return nil
}
