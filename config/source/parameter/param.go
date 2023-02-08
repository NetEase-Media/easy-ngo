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

package parameter

import (
	"os"
	"strings"
)

const (
	defaultPrefix = "-D" // default prefix for parameter
)

// param default can parse parameter with prefix "-D"
// such as -Dkey=value
// if you set the Prefix to "-C" it can also parse "-C" prefix parameter
// such as -Ckey=value
// but it just for inner package.
// will not allow use outside of the package.
type param struct {
	Prefix string
}

func New() *param {
	return &param{Prefix: defaultPrefix}
}

func (p *param) Load(sourcePathes []string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for _, arg := range os.Args {
		if len(arg) <= len(p.Prefix) || arg[:len(p.Prefix)] != p.Prefix {
			continue
		}
		kv := arg[len(p.Prefix):]
		kvs := strings.Split(kv, "=")
		if len(kvs) != 2 {
			continue
		}
		m[kvs[0]] = kvs[1]
	}
	return m, nil
}
