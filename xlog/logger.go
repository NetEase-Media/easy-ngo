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

package xlog

type Logger interface {
	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
	Warnf(format string, params ...interface{})
	Errorf(format string, params ...interface{})
	Panicf(format string, params ...interface{})
	DPanicf(format string, params ...interface{})
	Fatalf(format string, params ...interface{})
	Level() Level
}

var _ Logger = (*NopLogger)(nil)

func NewNopLogger() Logger {
	return &NopLogger{}
}

type NopLogger struct{}

func (n NopLogger) Debugf(format string, params ...interface{}) {
}

func (n NopLogger) Infof(format string, params ...interface{}) {
}

func (n NopLogger) Warnf(format string, params ...interface{}) {
}

func (n NopLogger) Errorf(format string, params ...interface{}) {
}

func (n NopLogger) Panicf(format string, params ...interface{}) {
}

func (n NopLogger) DPanicf(format string, params ...interface{}) {
}

func (n NopLogger) Fatalf(format string, params ...interface{}) {
}

func (n NopLogger) Level() Level {
	return FatalLevel
}
