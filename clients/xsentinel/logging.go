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

package xsentinel

import (
	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
	"github.com/alibaba/sentinel-golang/logging"
)

func NewLogger() *Logger {
	lg, _ := xfmt.Default()
	return &Logger{
		// logger: log.DefaultLogger().(*log.NgoLogger).WithOptions(zap.AddCallerSkip(1)),
		logger: lg,
	}
}

type Logger struct {
	logging.Logger
	logger xlog.Logger
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	// if !l.DebugEnabled() {
	// 	return
	// }
	l.logger.Debugf(msg, keysAndValues)
}

// func (l *Logger) DebugEnabled() bool {
// 	return l.logger.GetLevel() >= log.DebugLevel
// }

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	if !l.InfoEnabled() {
		return
	}
	l.logger.Infof(msg, keysAndValues)
}

//	func (l *Logger) InfoEnabled() bool {
//		return l.logger.GetLevel() >= log.InfoLevel
//	}
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	if !l.WarnEnabled() {
		return
	}
	l.logger.Warnf(msg, keysAndValues)
}

//	func (l *Logger) WarnEnabled() bool {
//		return l.logger.GetLevel() >= log.WarnLevel
//	}
func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	if !l.ErrorEnabled() {
		return
	}
	l.logger.Errorf(msg, keysAndValues)
}

// func (l *Logger) ErrorEnabled() bool {
// 	return l.logger.GetLevel() >= log.ErrorLevel
// }
