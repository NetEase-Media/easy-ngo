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

package rxxljob

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngots/xxxljob"
	conf "github.com/NetEase-Media/easy-ngog"
)

const (
	key = "ngo.client.xxxljob"
)

var xJobManager map[string]*xxxljob.XJobManager

func init() {
	hooks.Register(hooks.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	opt := xxxljob.Option{}
	err := conf.Get(key, &opt)
	if err != nil {
		panic(err)
	}
	cli := xxxljob.New(&opt, nil)
	k := key
	set(k, cli)
	return nil
}

func set(key string, client *xxxljob.XJobManager) {
	if xJobManager == nil {
		xJobManager = make(map[string]*xxxljob.XJobManager, 1)
	}
	xJobManager[key] = client
}

func GetXJobManager(key string) (cli *xxxljob.XJobManager) {
	var ok bool
	cli, ok = xJobManager[key]
	if !ok {
		return nil
	}
	return cli
}
