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
