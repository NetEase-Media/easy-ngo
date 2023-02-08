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
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
	"github.com/go-redis/redis/v8"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestCommands(t *testing.T) {
	do(func(redis Redis) {
		for i := 0; i < 10; i++ {
			_, err := redis.Del(context.Background(), "sharded:testkey1").Result()
			assert.NoError(t, err)
			result, err := redis.Set(context.Background(), "sharded:testkey"+strconv.Itoa(i), "testvalue"+strconv.Itoa(i), time.Second*5).Result()
			assert.NoError(t, err)
			assert.Equal(t, "OK", result)
			value, err := redis.Get(context.Background(), "sharded:testkey"+strconv.Itoa(i)).Result()
			assert.NoError(t, err)
			assert.Equal(t, "testvalue"+strconv.Itoa(i), value)
		}

		client := redis.(*ShardedClient)
		for _, c := range client.getAllShards() {
			container := c.client.(*RedisContainer)
			keys, err := container.Keys(context.Background(), "*").Result()
			fmt.Printf("[%v].Keys: %+v, err: %+v\n", container.Opt.Name, keys, err)
		}
	})
}

func TestPipeline(t *testing.T) {
	ctx := context.Background()
	do(func(client Redis) {
		client.Set(ctx, "k1", "v1", time.Second*60)
		client.Set(ctx, "k2", "v2", time.Second*60)
		client.Set(ctx, "k3", "v3", time.Second*60)
		client.Set(ctx, "k4", "v4", time.Second*60)
		client.Set(ctx, "k5", "v5", time.Second*60)

		pipe := client.Pipeline()
		pipe.Get(ctx, "k1")
		pipe.Get(ctx, "k2")
		pipe.Get(ctx, "k3")
		pipe.Get(ctx, "k4")
		pipe.Get(ctx, "k5")
		cmds, err := pipe.Exec(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "v1", cmds[0].(*redis.StringCmd).Val())
		assert.Equal(t, "v2", cmds[1].(*redis.StringCmd).Val())
		assert.Equal(t, "v3", cmds[2].(*redis.StringCmd).Val())
		assert.Equal(t, "v4", cmds[3].(*redis.StringCmd).Val())
		assert.Equal(t, "v5", cmds[4].(*redis.StringCmd).Val())

		c := client.(*ShardedClient)
		for _, cc := range c.getAllShards() {
			container := cc.client.(*RedisContainer)
			keys, err := container.Keys(context.Background(), "*").Result()
			fmt.Printf("[%v].Keys: %+v, err: %+v\n", container.Opt.Name, keys, err)
		}
	})
}

func TestPipelineUnsupportedCmd(t *testing.T) {
	ctx := context.Background()
	do(func(redis Redis) {
		pipe := redis.Pipeline()
		pipe.Dump(ctx, "key0")
		cmds, err := pipe.Exec(ctx)
		assert.Equal(t, "unsupport command: [dump]", err.Error())
		assert.Equal(t, 1, len(cmds))
	})
}

func BenchmarkShardedClient_Set(b *testing.B) {
	ctx := context.Background()
	do(func(redis Redis) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cmd := redis.Set(ctx, "test:"+strconv.Itoa(i), "value"+strconv.Itoa(i), time.Second*10)
			assert.NoError(b, cmd.Err())
		}
		b.StopTimer()
	})
}

func do(commandFn func(Redis)) {
	s1, _ := miniredis.Run()
	s2, _ := miniredis.Run()
	s3, _ := miniredis.Run()
	defer func() {
		s1.Close()
		s2.Close()
		s3.Close()
	}()
	opt1 := &Option{
		Name: fmt.Sprintf("test client %d", 1),
		Addr: []string{s1.Addr()},
	}
	rc1 := NewClient(opt1, &xfmt.XFmt{}, nil, nil)
	opt2 := &Option{
		Name: fmt.Sprintf("test client %d", 2),
		Addr: []string{s2.Addr()},
	}
	rc2 := NewClient(opt2, &xfmt.XFmt{}, nil, nil)
	opt3 := &Option{
		Name: fmt.Sprintf("test client %d", 3),
		Addr: []string{s3.Addr()},
	}
	rc3 := NewClient(opt3, &xfmt.XFmt{}, nil, nil)

	sis := []*ShardInfo{
		{
			id:     "shard-0",
			name:   "shard-0",
			client: rc1,
			weight: 1,
		},
		{
			id:     "shard-1",
			name:   "shard-1",
			client: rc2,
			weight: 1,
		},
		{
			id:     "shard-2",
			name:   "shard-2",
			client: rc3,
			weight: 1,
		},
	}
	client := NewShardedClient(sis)
	commandFn(client)
	client.Close()
}
