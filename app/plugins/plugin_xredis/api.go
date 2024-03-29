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

package pluginxredis

import (
	"sync"

	"github.com/NetEase-Media/easy-ngo/clients/xredis"
)

var (
	mu           sync.RWMutex
	redisClients = make(map[string]xredis.Redis)
)

func set(name string, client xredis.Redis) {
	mu.Lock()
	defer mu.Unlock()
	redisClients[name] = client
}

func GetClientByKey(name string) xredis.Redis {
	mu.RLock()
	defer mu.RUnlock()
	return redisClients[name]
}

func GetClient() xredis.Redis {
	return GetClientByKey("default")
}
