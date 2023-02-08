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

package xkafka

import (
	"testing"
	"time"

	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
	"github.com/stretchr/testify/assert"
)

const (
	NAME = "ngo-test"
)

func TestInit_InitProcess(t *testing.T) {
	opts := NewDefaultOptions()
	opts.Name = NAME
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	opts.Consumer.Group = "ngo"
	k, err := New(opts, &xfmt.XFmt{}, nil, nil)
	assert.Equal(t, nil, err)
	consumer := k.Consumer
	assert.Equal(t, "ngo", consumer.opt.Consumer.Group)
	producer := k.Producer
	assert.Equal(t, time.Second*10, producer.opt.Producer.Timeout)
}
