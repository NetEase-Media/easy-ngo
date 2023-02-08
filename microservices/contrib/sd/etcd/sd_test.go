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

package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/NetEase-Media/easy-ngo/microservices/sd"
	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func TestServiceDiscovery(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 3,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	assert.NoError(t, err)
	defer cli.Close()

	log, err := xfmt.Default()
	assert.NoError(t, err)

	s := New(ctx, cli, WithLogger(log))
	assert.Equal(t, s.opts.namespace, "microservices")

	info := &sd.ServiceInfo{Name: "test", Addr: "127.0.0.1"}

	watcher, err := s.NewWatcher(ctx, info.Name)
	assert.NoError(t, err)
	go func() {
		for updates := range watcher {
			t.Logf("watch: %d", len(updates))
			for _, update := range updates {
				t.Logf("next: %+v", update)
			}
		}
	}()

	err = s.Register(ctx, info)
	assert.NoError(t, err)

	infos, err := s.GetService(ctx, info.Name)
	assert.NoError(t, err)
	assert.NotEmpty(t, infos)
	assert.Equal(t, infos[0].Name, info.Name)

	err = s.Deregister(ctx, info)
	assert.NoError(t, err)

	infos, err = s.GetService(ctx, info.Name)
	assert.Empty(t, infos)

	time.Sleep(time.Second)
}

func TestKeepalive(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 3,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	assert.NoError(t, err)
	defer cli.Close()

	log, err := xfmt.Default()
	assert.NoError(t, err)

	s := New(ctx, cli, WithLogger(log))
	assert.Equal(t, s.opts.namespace, "microservices")

	info := &sd.ServiceInfo{Name: "test", Addr: "127.0.0.1"}
	err = s.Register(ctx, info)
	assert.NoError(t, err)

	lease := s.getLeaseID()
	s.client.Revoke(ctx, lease)

	assert.Eventually(t, func() bool {
		return s.getLeaseID() != 0 && lease != s.getLeaseID()
	}, time.Second*5, time.Second)

	ttl, err := s.client.TimeToLive(ctx, lease)
	assert.Nil(t, err)
	assert.Equal(t, int64(-1), ttl.TTL)

	ttl, err = s.client.TimeToLive(context.Background(), s.getLeaseID())
	assert.Nil(t, err)
	assert.True(t, ttl.TTL > 0)
}
