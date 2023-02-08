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
	"crypto/tls"
	"errors"
	"time"
)

// Options 是redis客户端的通用配置选项，兼容单实例和cluster类型。
type Option struct {
	// 用户需要保证名字唯一
	Name string

	// 连接类型，必须指定。包含client、cluster、sentinel、sharded_sentinel四种类型。
	ConnType string

	// 地址列表，格式为host:port。如果是单实例只会取第一个。
	Addr []string

	// master 名称，只当sentinel、sharded_sentinel 类型必填。如果是sentinel只会取第一个。
	MasterNames []string

	// 自动生成分片名称，如果为false，默认使用MasterName， 只当sharded_sentinel 类型使用。
	// 该字段用来兼容旧项目，非特殊情况请勿设置成true，否则在MasterNames顺序变化时会造成分配rehash
	AutoGenShardName bool

	// 用于认证的用户名
	Username string

	// 用于认证的密码
	Password string

	// 所使用的数据库
	DB int

	// 最大重试次数
	MaxRetries int

	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	// 超时时间
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// 最大连接数
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration

	// TODO: 未来增加
	TLSConfig    *tls.Config
	EnableTracer bool
}

func NewDefaultOptions() *Option {
	return &Option{}
}

func checkOptions(opt *Option) error {
	if opt.Name == "" {
		return errors.New("client name can not be nil")
	}
	if len(opt.Addr) == 0 {
		return errors.New("empty address")
	}
	return nil
}
