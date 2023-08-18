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

package xmemcache

import (
	"context"

	"github.com/bradfitz/gomemcache/memcache"
)

func New(config *Config) (*MemcacheProxy, error) {
	c := memcache.New(config.Addr...)
	c.Timeout = config.Timeout
	c.MaxIdleConns = config.MaxIdleConns
	p := &MemcacheProxy{
		base:  c,
		hooks: make([]Hook, 0),
	}
	return p, nil
}

type Hook interface {
	Before(context.Context, string, ...string) (context.Context, error)
	After(context.Context, error)
}

type MemcacheProxy struct {
	Config *Config
	base   *memcache.Client
	hooks  []Hook
}

func (mp *MemcacheProxy) Initialize() error {
	return nil
}

func (mp *MemcacheProxy) Destory() error {
	return nil
}

func (mp *MemcacheProxy) AddHook(h Hook) {
	mp.hooks = append(mp.hooks, h)
}
