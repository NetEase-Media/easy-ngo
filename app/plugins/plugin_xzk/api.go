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
