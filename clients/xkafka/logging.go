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

package xkafka

import (
	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/NetEase-Media/easy-ngoxfmt"
)

func NewLogger() *logger {
	dxfmt, _ := xfmt.Default()
	return &logger{
		logger: dxfmt, // TODO
	}
}

type logger struct {
	logger xlog.Logger
}

func (l *logger) Print(v ...interface{}) {
	for _, item := range v {
		l.logger.Infof("%v", item)
	}
}
func (l *logger) Printf(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}
func (l *logger) Println(v ...interface{}) {
	for _, item := range v {
		l.logger.Infof("%v", item)
	}
}
