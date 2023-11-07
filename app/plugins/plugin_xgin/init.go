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

package pluginxgin

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/NetEase-Media/easy-ngo/config"
	"github.com/NetEase-Media/easy-ngo/server/contrib/xgin"
)

func init() {
	app.RegisterPlugin(app.Initialize, Initialize)
	app.RegisterPlugin(app.Starting, Serve)
	app.RegisterPlugin(app.Stopping, Shutdown)
}

func Initialize(ctx context.Context) error {
	c := xgin.DefaultConfig()
	if err := config.UnmarshalKey("server", c); err != nil {
		return err
	}
	WithServer(xgin.New(c))
	return GetServer().Init()
}

func Serve(ctx context.Context) error {
	return GetServer().Serve()
}

func Shutdown(ctx context.Context) error {
	return GetServer().Shutdown()
}
