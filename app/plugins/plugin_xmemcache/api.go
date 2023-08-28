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
