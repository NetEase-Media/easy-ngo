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

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/NetEase-Media/easy-ngo/xlog"
)

type XZap struct {
	lv  *zap.AtomicLevel
	zl  *zap.Logger
	zsl *zap.SugaredLogger
}

func New(config *Config) (xlog.Logger, error) {
	//默认只记录info以上的日志
	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, err
	}
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        config.TimeKey,
		LevelKey:       config.LevelKey,
		NameKey:        config.NameKey,
		CallerKey:      config.CallerKey,
		FunctionKey:    config.FunctionKey,
		MessageKey:     config.MessageKey,
		StacktraceKey:  config.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(xlog.TIMEFORMAT),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	if !config.Console {
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	var encoder zapcore.Encoder
	switch config.Format {
	case JSON:
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}
	var core zapcore.Core
	if config.Console {
		core = zapcore.NewCore(encoder, os.Stdout, lv)
	} else {
		rlog, err := newRotateLog(config, config.Path, config.Suffix)
		if err != nil {
			return nil, err
		}
		relog, err := newRotateLog(config, config.ErrorPath, config.ErrorSuffix)
		if err != nil {
			return nil, err
		}
		elv := zapcore.ErrorLevel
		if errlogLevel, err := zapcore.ParseLevel(string(config.ErrlogLevel)); err == nil {
			elv = errlogLevel
		}
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(rlog), lv),
			zapcore.NewCore(encoder, zapcore.AddSync(relog), elv),
		)
	}
	zapOption := make([]zap.Option, 0)
	if config.WritableCaller {
		zapOption = append(zapOption, zap.AddCaller(), zap.AddCallerSkip(config.AddCallerSkip))
	}
	if config.WritableStack {
		zapOption = append(zapOption, zap.AddStacktrace(zapcore.ErrorLevel))
	}
	zl := zap.New(
		core,
		zapOption...,
	)
	return &XZap{
		lv:  &lv,
		zl:  zl,
		zsl: zl.Sugar(),
	}, nil
}

func newRotateLog(config *Config, p, suffix string) (io.Writer, error) {
	dir, err := filepath.Abs(p)
	if err != nil {
		return nil, err
	}
	linkName := path.Join(dir, config.FileName+"."+suffix)
	return &lumberjack.Logger{
		Filename:   linkName,
		MaxSize:    config.MaxSize, // megabytes
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge, // days
		LocalTime:  true,
		Compress:   config.Compress, // disabled by default
	}, nil
}

func (x *XZap) Debug(msg string, fields ...zap.Field) {
	x.zl.Debug(msg, fields...)
}

func (x *XZap) Info(msg string, fields ...zap.Field) {
	x.zl.Info(msg, fields...)
}

func (x *XZap) SugaredLogger() *zap.SugaredLogger {
	return x.zsl
}

func (x *XZap) Logger() *zap.Logger {
	return x.zl
}

func (x *XZap) Debugf(msg string, params ...interface{}) {
	x.zsl.Debugf(msg, params...)
}

func (x *XZap) Infof(msg string, params ...interface{}) {
	x.zsl.Infof(msg, params...)
}

func (x *XZap) Errorf(msg string, params ...interface{}) {
	x.zsl.Errorf(msg, params...)
}

func (x *XZap) Warnf(msg string, params ...interface{}) {
	x.zsl.Warnf(msg, params...)
}

func (x *XZap) Fatalf(msg string, params ...interface{}) {
	x.zsl.Fatalf(msg, params...)
}

func (x *XZap) Panicf(msg string, params ...interface{}) {
	x.zsl.Panicf(msg, params...)
}
