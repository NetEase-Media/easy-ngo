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

package xzk

import (
	"sync"

	"github.com/go-zookeeper/zk"
)

func New(config *Config) (*ZookeeperProxy, error) {
	proxy, err := newZookeeperProxy(config)
	if err != nil {
		return nil, err
	}
	return proxy, nil
}

type ZookeeperProxy struct {
	Config     *Config
	Conn       *zk.Conn
	tmpNode    sync.Map
	listenCh   chan Event
	listenSign uint32
	stop       chan struct{}
	done       *sync.WaitGroup
}
