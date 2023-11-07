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

package pluginxmemcache

import (
	"github.com/NetEase-Media/easy-ngo/clients/xmemcache"
)

var memecacheClients map[string]*xmemcache.MemcacheProxy

func set(key string, client *xmemcache.MemcacheProxy) {
	if memecacheClients == nil {
		memecacheClients = make(map[string]*xmemcache.MemcacheProxy, 1)
	}
	memecacheClients[key] = client
}

func GetClientByKey(key string) (cli *xmemcache.MemcacheProxy) {
	var ok bool
	cli, ok = memecacheClients[key]
	if !ok {
		return nil
	}
	return cli
}

func GetClient() (cli *xmemcache.MemcacheProxy) {
	var ok bool
	cli, ok = memecacheClients["default"]
	if !ok {
		return nil
	}
	return cli
}
