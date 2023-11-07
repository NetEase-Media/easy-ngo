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

package xtracer

type EXPORTER_NAME string

const (
	EXPORTER_NAME_JAEGER EXPORTER_NAME = "jaeger"
	EXPORTER_NAME_OLTP   EXPORTER_NAME = "oltp"
	EXPORTER_NAME_STDOUT EXPORTER_NAME = "stdout"
)

type Config struct {
	// 采样率
	SampleRate float64
	// 采样器
	ExporterName EXPORTER_NAME
	// OLTP采样器服务地址
	ExporterEndpoint string
	// OLTP采样器服务名称
	ServiceName string
}

func DefaultConfig() *Config {
	return &Config{
		SampleRate:   1,
		ExporterName: EXPORTER_NAME_STDOUT,
	}
}
