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

package rlog

import (
	"context"

	conf "github.com/NetEase-Media/easy-ngo/config"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngo/application/r/rlog"

	"github.com/NetEase-Media/easy-ngo/xlog/contrib/xzap"
)

const (
	key_xzap = rlog.Key_prefix + ".xzap"
)

func init() {
	hooks.Register(hooks.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	xzapOpts := make([]xzap.Option, 0)
	if err := conf.Get(key_xzap, &xzapOpts); err != nil {
		panic("load xzap config failed.")
	}
	if len(xzapOpts) == 0 {
		panic("no xzap config!")
	}
	for _, xzapOpt := range xzapOpts {
		xzap, err := xzap.New(&xzapOpt)
		if err != nil {
			panic("init xzap failed.")
		}
		rlog.Set(xzapOpt.Name, xzap)
	}
	return nil
}
