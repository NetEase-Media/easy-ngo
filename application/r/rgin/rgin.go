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

package rgin

import (
	"context"

	conf "github.com/NetEase-Media/easy-ngo/config"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngo/application/r/rmetrics"
	"github.com/NetEase-Media/easy-ngo/servers/xgin"
)

var server *xgin.Server

func init() {
	hooks.Register(hooks.Initialize, Initialize)
	hooks.Register(hooks.Start, Start)
	hooks.Register(hooks.Stop, Shutdown)
}

func Gin() *xgin.Server {
	return server
}

func Initialize(ctx context.Context) error {
	option := xgin.DefaultOption()
	conf.Get("ngo.server.gin", option)
	metrics := rmetrics.GetProvider()
	if server == nil {
		server = xgin.New(option, nil, metrics, nil)
	}
	server.Initialize()
	return nil
}

func Start(ctx context.Context) error {
	return server.Serve()
}

func Shutdown(ctx context.Context) error {
	return nil
}
