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
