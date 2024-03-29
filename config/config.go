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

package config

import "time"

type Config interface {
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetTime(key string) time.Time
	GetFloat64(key string) float64
	UnmarshalKey(key string, rawVal interface{}) error

	Init(protocols ...string) error
}
