package pluginxlog

import (
	logger "github.com/NetEase-Media/easy-ngo/xlog"
)

var xlog logger.Logger

func Debugf(msg string, params ...interface{}) {
	xlog.Debugf(msg, params...)
}

func Infof(msg string, params ...interface{}) {
	xlog.Infof(msg, params...)
}

func Warnf(msg string, params ...interface{}) {
	xlog.Warnf(msg, params...)
}

func Errorf(msg string, params ...interface{}) {
	xlog.Errorf(msg, params...)
}

func Fatalf(msg string, params ...interface{}) {
	xlog.Fatalf(msg, params...)
}

func Panicf(msg string, params ...interface{}) {
	xlog.Panicf(msg, params...)
}

func WithVendor(log logger.Logger) {
	xlog = log
}
