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
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
)

func TestMetricHook(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	w := newTestClientWrapper()
	defer w.Stop()

	c := &RedisContainer{
		Redis: w.client,
		Opt: Option{
			Name: "test",
			Addr: []string{w.server.Addr()},
		},
		redisType: RedisTypeClient,
	}
	client := w.client.Redis.(*redis.Client)
	client.AddHook(newMetricHook(c, nil, nil))
	c.Get(context.Background(), "a")
}
