package xlog

type Logger interface {
	Debugf(msg string, params ...interface{})
	Infof(msg string, params ...interface{})
}
