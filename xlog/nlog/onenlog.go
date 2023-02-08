package nlog

var _nlog *Nlog

func init() {
	_nlog, _ = New(DefaultOption())
}
func SetLevel(levelStr string) error {
	return _nlog.SetLevel(levelStr)
}
func GetLevel() string {
	return _nlog.GetLevel()
}
func SetFlags(flagStr string) {
	_nlog.SetFlags(flagStr)
}
func GetFlags() int {
	return _nlog.GetFlags()
}
func GetName() string {
	return _nlog.name
}
func String() string {
	return _nlog.String()
}
func Debugf(format string, fields ...interface{}) {
	_nlog.Debugf(format, fields...)
}
func Infof(format string, fields ...interface{}) {
	_nlog.Infof(format, fields...)
}
func Warnf(format string, fields ...interface{}) {
	_nlog.Warnf(format, fields...)
}
func Errorf(format string, fields ...interface{}) {
	_nlog.Errorf(format, fields...)
}
func DPanicf(format string, fields ...interface{}) {
	_nlog.DPanicf(format, fields...)
}
func Panicf(format string, fields ...interface{}) {
	_nlog.Panicf(format, fields...)
}
func Fatalf(format string, fields ...interface{}) {
	_nlog.Fatalf(format, fields...)
}
