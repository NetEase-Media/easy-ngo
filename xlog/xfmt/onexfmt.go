package xfmt

import "github.com/NetEase-Media/easy-ngo/xlog"

var _xfmt *XFmt

func init() {
	_xfmt, _ = New(DefaultOption())
}
func SetLevel(levelStr string) error {
	alevel, err := xlog.ParseLevel(levelStr)
	if err != nil {
		return err
	}
	_xfmt.level = alevel
	return nil
}
func GetLevel() string {
	return _xfmt.level.String()
}
func GetName() string {
	return _xfmt.name
}
func String() string {
	return "&{" + _xfmt.name + " " + _xfmt.level.String() + "}"
}
func Debugf(format string, fields ...interface{}) {
	_xfmt.Debugf(format, fields...)
}
func Infof(format string, fields ...interface{}) {
	_xfmt.Infof(format, fields...)
}
func Warnf(format string, fields ...interface{}) {
	_xfmt.Warnf(format, fields...)
}
func Errorf(format string, fields ...interface{}) {
	_xfmt.Errorf(format, fields...)
}
func DPanicf(format string, fields ...interface{}) {
	_xfmt.DPanicf(format, fields...)
}
func Panicf(format string, fields ...interface{}) {
	_xfmt.Panicf(format, fields...)
}
func Fatalf(format string, fields ...interface{}) {
	_xfmt.Fatalf(format, fields...)
}
