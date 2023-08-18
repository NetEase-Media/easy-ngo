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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 先使用kafka的zk测试
// var testZkClientAddr = "kafka"

var testZkClientAddr = "127.0.0.1:2181"

func TestInit(t *testing.T) {
	o := Config{
		Name:           "m",
		Addr:           []string{testZkClientAddr},
		SessionTimeout: time.Second * 5,
	}
	_, err := New(&o)
	assert.Equal(t, nil, err, "encounter error.")
}

func TestNewClientFromZookeeperException1(t *testing.T) {
	o := Config{
		Name:           "",
		Addr:           []string{testZkClientAddr},
		SessionTimeout: time.Second * 5,
	}
	// name=""
	_, err := New(&o)
	assert.Error(t, err)
}

func TestNewClientFromZookeeperException2(t *testing.T) {
	o := Config{
		Name:           "m",
		Addr:           []string{},
		SessionTimeout: time.Second * 5,
	}
	_, err := New(&o)
	assert.Error(t, err)
}

func TestNewClientFromZookeeperException3(t *testing.T) {
	o := Config{
		Name:           "m",
		Addr:           []string{"a"},
		SessionTimeout: time.Second * 1,
	}
	_, err := New(&o)
	assert.Error(t, err)
}

func TestOptionsProxy_Close(t *testing.T) {
	o := Config{
		Name: NAME,
		// Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 5,
	}
	c, err := New(&o)
	assert.Equal(t, nil, err, "")
	isE_00, re_00 := c.Exist("/pushTest")
	assert.Equal(t, nil, re_00, "判断存在出错")
	if !isE_00 {
		_, re_00 = c.CreateNode("/pushTest", PERSISTENT, "ddd")
		assert.Equal(t, nil, re_00, "创建节点出错")
	}
	time.Sleep(3 * time.Second)
	ch, err := c.StartLi()
	assert.Equal(t, nil, err, "状态监听已经开启")
	time.Sleep(1 * time.Second)
	c.Close()
	time.Sleep(2 * time.Second)
	event, ok := <-ch
	if !ok {
		fmt.Println("channel closed")
	}
	a, _, err := c.Conn.Exists("/pushTest/test3")
	assert.NotNil(t, err)
	assert.Equal(t, false, a, "")
	fmt.Println(event)
	assert.Equal(t, "StateDisconnected", c.GetConnState())
	if err != nil {
		_, bbb := c.CreateNode("/pushTest/test3", 1, "flag0")
		assert.NotNil(t, bbb)
	}
}
