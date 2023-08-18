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

package xxxljob

import (
	xxl "github.com/xxl-job/xxl-job-executor-go"
)

type XJobManager struct {
	config *Config
	xxl.Executor
}

func New(config *Config) *XJobManager {
	exec := xxl.NewExecutor(
		xxl.ServerAddr(config.Addr),
		xxl.AccessToken(config.Token),
		xxl.ExecutorIp(config.ExecutorIP),
		xxl.ExecutorPort(config.ExecutorPort),
		xxl.RegistryKey(config.ExecutorName),
	)
	return &XJobManager{
		config:   config,
		Executor: exec,
	}
}

func (m *XJobManager) Init() {
	m.Executor.Init()
}
