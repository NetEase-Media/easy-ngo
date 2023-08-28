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
