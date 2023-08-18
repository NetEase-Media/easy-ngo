package pluginxgrom

import "github.com/NetEase-Media/easy-ngo/clients/xgorm"

var dbClients map[string]*xgorm.Client

func set(key string, client *xgorm.Client) {
	if dbClients == nil {
		dbClients = make(map[string]*xgorm.Client, 1)
	}
	dbClients[key] = client
}

func GetDBClientByKey(key string) (cli *xgorm.Client) {
	var ok bool
	cli, ok = dbClients[key]
	if !ok {
		return nil
	}
	return cli
}

func GetDBClient() (cli *xgorm.Client) {
	return GetDBClientByKey("default")
}
