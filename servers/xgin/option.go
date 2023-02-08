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

package xgin

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type Option struct {
	Host           string
	Port           int
	Mode           string
	EnabledMetric  bool
	EnableTracer   bool
	ServiceAddress string
	MetricsPath    string

	ShutdownTimeout time.Duration
	Middlewares     *MiddlewaresOption
}

type MiddlewaresOption struct {
	AccessLog *AccessLogMwOption
	// JwtAuth   *JwtAuthMwOption
}

func DefaultOption() *Option {
	return &Option{
		Host:          "0.0.0.0",
		Port:          8080,
		Mode:          gin.DebugMode,
		EnabledMetric: false,
		EnableTracer:  false,
		MetricsPath:   "/metrics",
		Middlewares: &MiddlewaresOption{
			AccessLog: &AccessLogMwOption{
				Enabled: false,
			},
		},
	}
}

func (option *Option) Address() string {
	return fmt.Sprintf("%s:%d", option.Host, option.Port)
}
