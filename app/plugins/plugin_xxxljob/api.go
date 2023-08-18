package pluginxxxljob

import (
	"github.com/NetEase-Media/easy-ngo/clients/xxxljob"
)

var xJobManager map[string]*xxxljob.XJobManager

func set(key string, client *xxxljob.XJobManager) {
	if xJobManager == nil {
		xJobManager = make(map[string]*xxxljob.XJobManager, 1)
	}
	xJobManager[key] = client
}

func GetXJobManagerByKey(key string) (cli *xxxljob.XJobManager) {
	var ok bool
	cli, ok = xJobManager[key]
	if !ok {
		return nil
	}
	return cli
}

func GetXJobManager() (cli *xxxljob.XJobManager) {
	var ok bool
	cli, ok = xJobManager["default"]
	if !ok {
		return nil
	}
	return cli
}
