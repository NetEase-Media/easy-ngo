package xlog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func New(conifg *Config) Logger {
	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	zapOption := make([]zap.Option, 0)
	var encoder zapcore.Encoder
	core := zapcore.NewCore(encoder, os.Stdout, lv)
	zl := zap.New(
		core,
		zapOption...,
	)
	return &XZap{
		lv:  &lv,
		zl:  zl,
		zsl: zl.Sugar(),
	}
}

func (x *XZap) Debugf(msg string, params ...interface{}) {
	x.zsl.Debugf(msg, params...)
}

func (x *XZap) Infof(msg string, params ...interface{}) {
	x.zsl.Infof(msg, params...)
}
