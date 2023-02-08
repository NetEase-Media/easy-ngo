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

package rprometheus

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/application/hooks"
	"github.com/NetEase-Media/easy-ngocation/r/rmetrics"
	conf "github.com/NetEase-Media/easy-ngog"
	"github.com/NetEase-Media/easy-ngovability/contrib/xprometheus"
)

const (
	defaultNamespace = "ngo"
	defaultSubsystem = "app"
)

func init() {
	hooks.Register(hooks.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	namespace := conf.GetString("ngo.prometheus.namespace")
	subsystem := conf.GetString("ngo.prometheus.subsystem")
	if namespace == "" {
		namespace = defaultNamespace
	}
	if subsystem == "" {
		subsystem = defaultSubsystem
	}
	rmetrics.SetProvider(xprometheus.NewProvider(namespace, subsystem))
	return nil
}
