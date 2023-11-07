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

package pluginxzk

import "github.com/NetEase-Media/easy-ngo/clients/xzk"

var zkClients map[string]*xzk.ZookeeperProxy

func set(key string, client *xzk.ZookeeperProxy) {
	if zkClients == nil {
		zkClients = make(map[string]*xzk.ZookeeperProxy, 1)
	}
	zkClients[key] = client
}

func GetZKClientByKey(key string) (cli *xzk.ZookeeperProxy) {
	var ok bool
	cli, ok = zkClients[key]
	if !ok {
		return nil
	}
	return cli
}

func GetZKClient() (cli *xzk.ZookeeperProxy) {
	var ok bool
	cli, ok = zkClients["default"]
	if !ok {
		return nil
	}
	return cli
}
