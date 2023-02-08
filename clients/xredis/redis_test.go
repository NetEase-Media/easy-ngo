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
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	opts := &Option{
		Name:     "client1",
		Addr:     []string{"127.0.0.1:2379"},
		ConnType: "client",
	}

	client, err := newWithOption(opts, &xfmt.XFmt{}, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, []string{"127.0.0.1:2379"}, client.Opt.Addr)
	client.Close()
}

func generateKey() string {
	t := time.Now().Unix()
	return "test_key:" + strconv.FormatInt(t, 10) + "_" + strconv.FormatInt(int64(rand.Int()), 10)
}
