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
	"testing"

	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
	"github.com/stretchr/testify/assert"
)

func TestCRUD(t *testing.T) {
	o := Option{
		Name: "m1",
		Addr: []string{"127.0.0.1:11312"},
	}

	c, err := New(&o, &xfmt.XFmt{}, nil, nil)
	assert.Equal(t, nil, err, "encounter error.")
	assert.NotEqual(t, nil, c, "Init Client Failed")

	data := make(map[string]string)
	data["Halo"] = "World"
	data["Halo1"] = "World1"

	// 插入
	for k, v := range data {
		err := c.Set(context.Background(), k, v)
		assert.Equal(t, nil, err, "err.")
	}

	// 查询
	for k := range data {
		d, err := c.Get(context.Background(), k)
		expect := data[k]
		assert.Equal(t, nil, err, "err.")
		assert.Equal(t, expect, d, "wrong data.")
	}

	// 批量查询
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	m, err := c.MGet(context.Background(), keys)
	assert.Equal(t, nil, err, "err.")
	assert.Equal(t, data, m, "mget wrong data.")

	// 删除
	for k := range data {
		err := c.Delete(context.Background(), k)
		assert.Equal(t, nil, err, "err.")
	}

	// 校验删除结果
	m, err = c.MGet(context.Background(), keys)
	assert.Equal(t, nil, err, "err.")
	assert.Equal(t, 0, len(m), "delete failed")
}
