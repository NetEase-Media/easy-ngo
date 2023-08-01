package xlog

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	formatTXT   = "txt"
	formatJSON  = "json"
	formatBlank = "blank"
	timeFormat  = "2006-01-02 15:04:05.999"
)

type XZap struct {
	lv  *zap.AtomicLevel
	zl  *zap.Logger
	zsl *zap.SugaredLogger
}

func New(config *Config) (Logger, error) {
	//默认只记录info以上的日志
	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(config.Level)); err != nil {
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
	if config.Environment == "development" {
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	}
	var encoder zapcore.Encoder
	switch config.Format {
	case formatJSON:
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = NewJSONEncoder(config, encoderCfg)
	case formatTXT:
		encoder = NewConsoleEncoder(config, encoderCfg)
	case formatBlank:
		encoderCfg.TimeKey = ""
		encoderCfg.LevelKey = ""
		encoderCfg.NameKey = ""
		encoderCfg.CallerKey = ""
		encoderCfg.FunctionKey = ""
		encoderCfg.StacktraceKey = ""
		encoder = NewConsoleEncoder(config, encoderCfg)
	default:
		encoder = NewConsoleEncoder(config, encoderCfg)
	}
	var core zapcore.Core
	if config.Environment == "development" {
		core = zapcore.NewCore(encoder, os.Stdout, lv)
	} else {
		rlog, err := newRotateLog(config, config.Path, "log")
		if err != nil {
			return nil, err
		}
		relog, err := newRotateLog(config, config.ErrorPath, "error.log")
		if err != nil {
			return nil, err
		}

		elv := zapcore.ErrorLevel
		// if errlogLevel, err := zapcore.ParseLevel(config.ErrlogLevel); err == nil {
		// 	elv = errlogLevel
		// }

		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(rlog), lv),
			zapcore.NewCore(encoder, zapcore.AddSync(relog), elv),
		)
	}
	zapOption := make([]zap.Option, 0)
	// if config.WritableCaller {
	// 	zapOption = append(zapOption, zap.AddCaller(), zap.AddCallerSkip(opt.Skip))
	// }
	// if config.WritableStack {
	// 	zapOption = append(zapOption, zap.AddStacktrace(zapcore.ErrorLevel))
	// }
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

func (x *XZap) Debugf(msg string, params ...interface{}) {
	x.zsl.Debugf(msg, params...)
}

func (x *XZap) Infof(msg string, params ...interface{}) {
	x.zsl.Infof(msg, params...)
}
