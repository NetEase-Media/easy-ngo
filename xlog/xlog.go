package xlog

type Logger interface {
	Debugf(msg string, params ...interface{})
	Infof(msg string, params ...interface{})
	Errorf(msg string, params ...interface{})
	Warnf(msg string, params ...interface{})
	Fatalf(msg string, params ...interface{})
	Panicf(msg string, params ...interface{})
}
