# xlog 日志系统

## 简介
`easy-ngo`在新版本去掉了对Golang的fmt和log做的简单封装，只引入优秀的开源日志组件zap，对zap进行了方便使用的封装，可以让开发者方便的使用日志功能。开发者可以创建自定义配置的日志实例。
## 代码结构
### easy-ngo的日志位于xlog目录下，抽象日志接口Logger，提供了xzap实现
```golang
type Logger interface {
	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
	Warnf(format string, params ...interface{})
	Errorf(format string, params ...interface{})
	Panicf(format string, params ...interface{})
	Fatalf(format string, params ...interface{})
}
```
## 核心概念
### Level 日志级别，默认级别为Info
* Debug
* Info
* Warn
* Error
* Panic
* Fatal
### Option 是一个配置选项，可以用来配置日志的输出格式，输出位置，日志级别等。
* xzap的Option
  * Name 日志实例名
  * NoFile 是否为开发模式，如果是true，只显示到标准输出
  * Format 日志格式txt、json、blank，分别为文本、json、空格分割的格式
  * WritableStack 是否需要打印error及以上级别的堆栈信息
  * Skip 跳过的日志层数
  * WritableCaller 是否需要打印行号函数信息
  * Level 日志级别
  * Path 日志存储路径
  * FileName 日志文件名称
  * PackageLevel    map[string]string // 包级别日志等级设置
  * ErrlogLevel 错误日志级别，默认error
  * ErrorPath 错误日志存储路径
  * MaxAge 保留日志文件的最大天数，默认7天
  * MaxBackups 保留日志文件的最大个数，默认7个
  * MaxSize 在进行切割之前，日志文件的最大大小（以MB为单位）默认1024
  * Compress 是否压缩/归档旧文件
  * packageLogLevel map[string]zapcore.Level