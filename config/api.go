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

import (
	"time"
)

var config Config

func Exists(key string) bool {
	return Get(key) != nil
}

func Get(key string) interface{} {
	return config.Get(key)
}

func GetString(key string) string {
	return config.GetString(key)
}

func GetInt(key string) int {
	return config.GetInt(key)
}

func GetBool(key string) bool {
	return config.GetBool(key)
}

func GetTime(key string) time.Time {
	return config.GetTime(key)
}

func GetFloat64(key string) float64 {
	return config.GetFloat64(key)
}

func UnmarshalKey(key string, rawVal interface{}) error {
	return config.UnmarshalKey(key, &rawVal)
}

func WithConfig(c Config) {
	config = c
}
