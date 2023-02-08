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
	"github.com/NetEase-Media/easy-ngo/xlog"
	xxl "github.com/xxl-job/xxl-job-executor-go"
)

type XJobManager struct {
	option *Option
	xxl.Executor
}

func New(option *Option, log xlog.Logger) *XJobManager {
	exec := xxl.NewExecutor(
		xxl.ServerAddr(option.Addr),
		xxl.AccessToken(option.Token),
		xxl.ExecutorIp(option.ExecutorIP),
		xxl.ExecutorPort(option.ExecutorPort),
		xxl.RegistryKey(option.ExecutorName),
	)
	ret := &XJobManager{
		option:   option,
		Executor: exec,
	}
	ret.Init()
	// err := ret.Run()
	// if err != nil {
	// 	panic("xxxljob run error.")
	// }
	return ret
}
