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

package xzap

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/NetEase-Media/easy-ngo/xlog"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	formatTXT   = "txt"
	formatJSON  = "json"
	formatBlank = "blank"

	timeFormat        = "2006-01-02 15:04:05.999"
	DefaultLoggerName = "default"
)

var _ xlog.Logger = (*XZap)(nil)

type XZap struct {
	lv  *zap.AtomicLevel
	zl  *zap.Logger
	zsl *zap.SugaredLogger
	opt *Option
}

func Default() *XZap {
	log, _ := zap.NewDevelopment()
	return &XZap{zsl: log.Sugar()}
}

func New(opt *Option) (*XZap, error) {
	if err := checkOption(opt); err != nil {
		return nil, err
	}

	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(opt.Level)); err != nil {
		return nil, err
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "file",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(timeFormat),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if !opt.NoFile {
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	var encoder zapcore.Encoder

	switch opt.Format {
	case formatJSON:
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = NewJSONEncoder(opt, encoderCfg)
	case formatTXT:
		encoder = NewConsoleEncoder(opt, encoderCfg)
	case formatBlank:
		encoderCfg.TimeKey = ""
		encoderCfg.LevelKey = ""
		encoderCfg.NameKey = ""
		encoderCfg.CallerKey = ""
		encoderCfg.FunctionKey = ""
		encoderCfg.StacktraceKey = ""
		encoder = NewConsoleEncoder(opt, encoderCfg)
	default:
		encoder = NewConsoleEncoder(opt, encoderCfg)
	}

	var core zapcore.Core
	if opt.NoFile {
		core = zapcore.NewCore(encoder, os.Stdout, lv)
	} else {
		rlog, err := newRotateLog(opt, opt.Path, "log")
		if err != nil {
			return nil, err
		}
		relog, err := newRotateLog(opt, opt.ErrorPath, "error.log")
		if err != nil {
			return nil, err
		}

		elv := zapcore.ErrorLevel
		if errlogLevel, err := zapcore.ParseLevel(opt.ErrlogLevel); err == nil {
			elv = errlogLevel
		}

		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(rlog), lv),
			zapcore.NewCore(encoder, zapcore.AddSync(relog), elv),
		)
	}

	zapOption := make([]zap.Option, 0)
	if opt.WritableCaller {
		zapOption = append(zapOption, zap.AddCaller(), zap.AddCallerSkip(opt.Skip))
	}
	if opt.WritableStack {
		zapOption = append(zapOption, zap.AddStacktrace(zapcore.ErrorLevel))
	}
	// zapOption = append(zapOption, zap.Hooks(metricsHook()))

	zl := zap.New(
		core,
		zapOption...,
	)

	return &XZap{
		lv:  &lv,
		zl:  zl,
		zsl: zl.Sugar(),
		opt: opt,
	}, nil
}

// func metricsHook() func(zapcore.Entry) error {
// 	return func(entry zapcore.Entry) error {
// 		if entry.Level < zapcore.ErrorLevel {
// 			return nil
// 		}
// 		if !metrics.IsMetricsEnabled() {
// 			return nil
// 		}
// 		errType := "NONE"
// 		collectors.ExceptionCollector().RecordError(
// 			errType,
// 			entry.Caller.String(),
// 			entry.Message,
// 			entry.Stack,
// 		)
// 		return nil
// 	}
// }

func newRotateLog(opt *Option, p, suffix string) (io.Writer, error) {
	dir, err := filepath.Abs(p)
	if err != nil {
		return nil, err
	}

	linkName := path.Join(dir, opt.FileName+"."+suffix)

	return &lumberjack.Logger{
		Filename:   linkName,
		MaxSize:    opt.MaxSize, // megabytes
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge, // days
		LocalTime:  true,
		Compress:   opt.Compress, // disabled by default
	}, nil
}
func (xZap *XZap) Infof(format string, fields ...interface{}) {
	xZap.zsl.Infof(format, fields...)
}
func (xZap *XZap) Warnf(format string, fields ...interface{}) {
	xZap.zsl.Warnf(format, fields...)
}
func (xZap *XZap) Errorf(format string, fields ...interface{}) {
	xZap.zsl.Errorf(format, fields...)
}
func (xZap *XZap) Debugf(format string, fields ...interface{}) {
	xZap.zsl.Debugf(format, fields...)
}
func (xZap *XZap) Fatalf(format string, fields ...interface{}) {
	xZap.zsl.Fatalf(format, fields...)
}

func (xZap *XZap) DPanicf(format string, params ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (xZap *XZap) Panicf(format string, fields ...interface{}) {
	xZap.zsl.Panicf(format, fields...)
}

func (xZap *XZap) Level() xlog.Level {
	//TODO implement me
	panic("implement me")
}
