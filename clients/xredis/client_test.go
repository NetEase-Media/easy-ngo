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
	"strconv"
	"testing"
	"time"

	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

type testClientWrapper struct {
	server *miniredis.Miniredis
	client *RedisContainer
}

func newTestClientWrapper() *testClientWrapper {
	var w testClientWrapper
	var err error
	w.server, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	opt := &Option{
		Addr: []string{w.server.Addr()},
		Name: "test client",
	}
	w.client = NewClient(opt, &xfmt.XFmt{}, nil, nil)
	return &w
}

func (w *testClientWrapper) Stop() {
	if err := w.client.Close(); err != nil {
		panic(err)
	}
	w.server.Close()
}

func TestClient_Set(t *testing.T) {
	w := newTestClientWrapper()
	defer w.Stop()

	k1 := generateKey()
	str, _ := w.client.SetEX(context.Background(), k1, "1000", time.Second*20).Result()
	if str != "OK" {
		t.Fatal("SetEX not valid", str)
	}
	str2, _ := w.client.Get(context.Background(), k1).Result()
	if str2 != "1000" {
		t.Fatal("Get not valid", str2)
	}
}

func BenchmarkClient_Set(b *testing.B) {
	ctx := context.Background()
	w := newTestClientWrapper()
	defer w.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := w.client.Set(ctx, "test:"+strconv.Itoa(i), "value"+strconv.Itoa(i), time.Second*10)
		assert.NoError(b, cmd.Err())
	}
	b.StopTimer()
}

func TestClient_Del(t *testing.T) {
	w := newTestClientWrapper()
	defer w.Stop()

	k1 := generateKey()
	k2 := generateKey()
	w.client.SetEX(context.Background(), k1, "1000", time.Second*60).Result()
	w.client.SetEX(context.Background(), k2, "2000", time.Second*60).Result()
	keys := []string{k1, k2, k1 + "_temp1", k1 + "_temp2", k2 + "_temp"}
	delNum, err := w.client.Del(context.Background(), keys...).Result()
	if delNum != 2 {
		t.Fatal("Del not valid", delNum, err)
	}
}

// TODO Exists 多个key会报错？
func TestClient_Exists(t *testing.T) {
	w := newTestClientWrapper()
	defer w.Stop()

	k1 := generateKey()

	if num, err := w.client.Exists(context.Background(), k1).Result(); num != 0 {
		t.Fatal("got not valid", err)
	}
	w.client.SetEX(context.Background(), k1, "1", time.Second*60).Result()
	if num, err := w.client.Exists(context.Background(), k1).Result(); num != 1 {
		t.Fatal("got not valid", err)
	}
}

func TestClient_ZRange(t *testing.T) {
	w := newTestClientWrapper()
	defer w.Stop()

	k1 := generateKey()
	w.client.ZAdd(context.Background(), k1, &redis.Z{Score: 1.55, Member: "key1"}).Result()
	w.client.ZAdd(context.Background(), k1, &redis.Z{Score: 1.56, Member: "key2"}).Result()
	w.client.Expire(context.Background(), k1, time.Minute)

	vals, err := w.client.ZRange(context.Background(), k1, 0, -1).Result()
	if len(vals) != 2 {
		t.Fatal("got not valid", err)
	}
	if delNum, _ := w.client.Del(context.Background(), k1).Result(); delNum < 1 {
		t.Fatal("del not valid", err)
	}
}
