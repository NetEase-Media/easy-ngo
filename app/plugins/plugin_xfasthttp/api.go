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

package pluginxfasthttp

import (
	"github.com/NetEase-Media/easy-ngo/clients/xfasthttp"
)

var httpClients map[string]*xfasthttp.Xfasthttp

func set(key string, client *xfasthttp.Xfasthttp) {
	if httpClients == nil {
		httpClients = make(map[string]*xfasthttp.Xfasthttp, 1)
	}
	httpClients[key] = client
}

func GetXfasthttpByKey(key string) (cli *xfasthttp.Xfasthttp) {
	var ok bool
	cli, ok = httpClients[key]
	if !ok {
		return nil
	}
	return cli
}

func GetXfasthttp() (cli *xfasthttp.Xfasthttp) {
	var ok bool
	cli, ok = httpClients["default"]
	if !ok {
		return nil
	}
	return cli
}
