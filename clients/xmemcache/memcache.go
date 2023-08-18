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
	"os"

	"github.com/bradfitz/gomemcache/memcache"
)

func (m *MemcacheProxy) processBefore(ctx context.Context, operation string, key ...string) (context.Context, error) {
	var err error
	for i := 0; i < len(m.hooks); i++ {
		ctx, err = m.hooks[i].Before(ctx, operation, key...)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

func (m *MemcacheProxy) processAfter(ctx context.Context, err error) {
	for i := len(m.hooks) - 1; i > 0; i-- {
		m.hooks[i].After(ctx, err)
	}
}

// Get 根据Key获取缓存数值
func (m *MemcacheProxy) Get(ctx context.Context, key string) (string, error) {
	//缓存是否命中。针对修改类操作，如set，hits总是为true
	// var hits bool
	var item *memcache.Item
	var err error

	ctx, err = m.processBefore(ctx, "Get", key)
	defer m.processAfter(ctx, err)
	if err != nil {
		return "", err
	}

	item, err = m.base.Get(key)
	if err != nil {
		return "", err
	}
	// hits = item != nil && len(item.Value) > 0
	return string(item.Value), err
}

// MGet 获取多个数值
func (m *MemcacheProxy) MGet(ctx context.Context, keys []string) (map[string]string, error) {
	//缓存是否命中。针对修改类操作，如set，hits总是为true
	var hits bool
	var rets map[string]*memcache.Item
	var err error
	ctx, err = m.processBefore(ctx, "MGet", keys...)
	defer m.processAfter(ctx, err)
	if err != nil {
		return nil, err
	}

	rets, err = m.base.GetMulti(keys)

	if err != nil {
		return nil, err
	}

	if len(rets) == 0 {
		return nil, nil
	}

	r := make(map[string]string, len(rets))
	for _, v := range rets {
		if !hits {
			hits = v != nil && len(v.Value) > 0
		}
		r[v.Key] = string(v.Value)
	}
	return r, nil
}

// Set 设置缓存
func (m *MemcacheProxy) Set(ctx context.Context, key string, value string) error {
	var err error

	ctx, err = m.processBefore(ctx, "Set", key, value)
	defer m.processAfter(ctx, err)
	if err != nil {
		return err
	}

	item := memcache.Item{
		Key:   key,
		Value: []byte(value),
	}
	return m.base.Set(&item)
}

// SetWithExpire 设置缓存，并且添加超时
// expire 以s为单位
func (m *MemcacheProxy) SetWithExpire(ctx context.Context, key string, value string, expire int) error {
	var err error

	ctx, err = m.processBefore(ctx, "SetWithExpire", key, value)
	defer m.processAfter(ctx, err)
	if err != nil {
		return err
	}

	item := memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: int32(expire),
	}
	return m.base.Set(&item)

}

// Delete 删除操作
func (m *MemcacheProxy) Delete(ctx context.Context, key string) error {
	var err error
	ctx, err = m.processBefore(ctx, "Delete", key)
	defer m.processAfter(ctx, err)
	if err != nil {
		return err
	}
	return m.base.Delete(key)
}

func hostName() string {
	name, err := os.Hostname()
	if err != nil {
		return "default"
	}
	return name
}
