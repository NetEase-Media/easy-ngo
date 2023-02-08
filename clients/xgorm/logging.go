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

package xgorm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NetEase-Media/easy-ngo/xlog"
	"gorm.io/gorm/logger"
)

var (
	traceStr     = "[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s [%.3fms] [rows:%v] %s"
	traceErrStr  = "%s [%.3fms] [rows:%v] %s"
)

type xlogger struct {
	logger.Interface
	logger.Config
	logger xlog.Logger
}

func NewLogger(config logger.Config, logger xlog.Logger) *xlogger {
	return &xlogger{
		Config: config,
		logger: logger,
	}
}

// LogMode log mode
func (l *xlogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info print info
func (l xlogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if strings.HasSuffix(msg, "\n") {
		msg = msg[:len(msg)-1]
	}
	l.logger.Infof(msg, data...)
}

// Warn print warn messages
func (l xlogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if strings.HasSuffix(msg, "\n") {
		msg = msg[:len(msg)-1]
	}
	l.logger.Warnf(msg, data...)
}

// Error print error messages
func (l xlogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if strings.HasSuffix(msg, "\n") {
		msg = msg[:len(msg)-1]
	}
	l.logger.Errorf(msg, data...)
}

// Trace print sql message
func (l xlogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logger.Level() >= xlog.ErrorLevel:
		sql, rows := fc()
		if rows == -1 {
			l.logger.Errorf(traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.logger.Errorf(traceErrStr, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.logger.Level() >= xlog.WarnLevel:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.logger.Warnf(traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.logger.Warnf(traceWarnStr, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.logger.Level() >= xlog.InfoLevel:
		sql, rows := fc()
		if rows == -1 {
			l.logger.Infof(traceStr, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.logger.Infof(traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
