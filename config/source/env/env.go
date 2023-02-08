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

package env

import (
	"os"
	"strings"
)

const DefaultEnvPrefix = "APP_"

type Env struct {
	Keys   []string
	Prefix string
}

func New() *Env {
	env := &Env{Prefix: DefaultEnvPrefix}
	return env
}

func (env *Env) AddKeys(keys ...string) *Env {
	env.Keys = append(env.Keys, keys...)
	return env
}

func (env *Env) Load(sourcePathes []string) (map[string]interface{}, error) {
	rmap := make(map[string]interface{})
	if len(env.Keys) > 0 {
		for _, key := range env.Keys {
			rmap[key] = os.Getenv(key)
		}
	}
	environmentString := os.Environ()
	for _, envString := range environmentString {
		kv := strings.Split(envString, "=")
		if len(kv) != 2 {
			continue
		}
		if env.Prefix != "" {
			if len(kv[0]) >= len(env.Prefix) && kv[0][:len(env.Prefix)] == env.Prefix {
				rmap[kv[0]] = kv[1]
			}
		} else {
			rmap[kv[0]] = kv[1]
		}
	}
	return rmap, nil
}
