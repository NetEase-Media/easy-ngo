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

package xgorm

import (
	"time"
)

// mysql client configuration
type Option struct {
	Name            string
	Type            string
	Url             string
	MaxIdleCons     int
	MaxOpenCons     int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	EnableTracer    bool
}

func defaultOption() *Option {
	return &Option{
		Type:            "mysql",
		MaxIdleCons:     10,
		MaxOpenCons:     10,
		ConnMaxLifetime: time.Second * 1000,
		ConnMaxIdleTime: time.Second * 60,
	}
}

// func NewOption(key string) (*Option, error) {
// 	opt := defaultOption()
// 	if err := conf.Get(key, opt); err != nil {
// 		return nil, err
// 	}
// 	return opt, nil
// }
