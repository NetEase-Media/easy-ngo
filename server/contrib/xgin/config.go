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

import "github.com/NetEase-Media/easy-ngo/xmetrics"

type MODE string

const (
	DEBUG   MODE = "debug"
	RELEASE      = "release"
	TEST         = "test"
)

type Config struct {
	Host           string
	Port           int
	EnabledMetrics bool
	EnabledTracer  bool
	Mode           MODE
	Metrics        Metrics
}

type Metrics struct {
	Bucket           xmetrics.Bucket
	ExcludeByPrefix  []string
	ExcludeByRegular []string
	IncludeByPrefix  []string
	IncludeByRegular []string
}

func DefaultConfig() *Config {
	return &Config{
		Host:           "0.0.0.0",
		Port:           8080,
		EnabledMetrics: false,
		EnabledTracer:  false,
		Mode:           DEBUG,
	}
}
