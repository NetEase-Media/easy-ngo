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

package xfmt

import (
	"fmt"

	"github.com/NetEase-Media/easy-ngo/xlog"
)

type XFmt struct {
	name  string
	level xlog.Level
}

func Default() (*XFmt, error) {
	return New(DefaultOption())
}

func New(opt *Option) (*XFmt, error) {
	alevel, err := xlog.ParseLevel(opt.Level)
	if err != nil {
		return nil, err
	}
	return &XFmt{name: opt.Name, level: alevel}, nil
}

func (xFmt *XFmt) Debugf(format string, fields ...interface{}) {
	if xFmt.level <= xlog.DebugLevel {
		fmt.Printf("[\x1b[0;35m"+xlog.Level.CapitalString(xlog.DebugLevel)+"\x1b[0m] "+format+"\n", fields...)
	}
}

func (xFmt *XFmt) Infof(format string, fields ...interface{}) {
	if xFmt.level <= xlog.InfoLevel {
		fmt.Printf("[\x1b[0;34m"+xlog.Level.CapitalString(xlog.InfoLevel)+"\x1b[0m] "+format+"\n", fields...)
	}
}

func (xFmt *XFmt) Warnf(format string, fields ...interface{}) {
	if xFmt.level <= xlog.WarnLevel {
		fmt.Printf("[\x1b[0;33m"+xlog.Level.CapitalString(xlog.WarnLevel)+"\x1b[0m] "+format+"\n", fields...)
	}
}

func (xFmt *XFmt) Errorf(format string, fields ...interface{}) {
	if xFmt.level <= xlog.ErrorLevel {
		fmt.Printf("[\x1b[0;31m"+xlog.Level.CapitalString(xlog.ErrorLevel)+"\x1b[0m] "+format+"\n", fields...)
	}
}

func (xFmt *XFmt) DPanicf(format string, fields ...interface{}) {
	if xFmt.level <= xlog.DPanicLevel {
		fmt.Printf("[\x1b[0;31m"+xlog.Level.CapitalString(xlog.DPanicLevel)+"\x1b[0m] "+format+"\n", fields...)
	}
}

func (xFmt *XFmt) Panicf(format string, fields ...interface{}) {
	if xFmt.level <= xlog.PanicLevel {
		fmt.Printf("[\x1b[0;31m"+xlog.Level.CapitalString(xlog.PanicLevel)+"\x1b[0m] "+format+"\n", fields...)
	}
}

func (xFmt *XFmt) Fatalf(format string, fields ...interface{}) {
	if xFmt.level <= xlog.FatalLevel {
		fmt.Printf("[\x1b[0;31m"+xlog.Level.CapitalString(xlog.FatalLevel)+"\x1b[0m] "+format+"\n", fields...)
	}
}

func (xFmt *XFmt) Level() xlog.Level {
	return xFmt.level
}
