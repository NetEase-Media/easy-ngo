# xlog 日志系统

## 简介
easy-ngo对Golang的fmt和log做了简单封装，同时引入优秀的开源日志组件zap，对zap进行了方便使用的封装，可以让开发者方便的使用日志功能。开发者可以创建自定义配置的日志实例。同时分别在xfmt和nlog包中，提供了默认的日志实例，开发者可以直接使用，无需配置。
## 代码结构
### easy-ngo的日志位于xlog目录下，抽象日志接口Logger，提供了xfmt,nlog,xzap三种实现
```golang
type Logger interface {
	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
	Warnf(format string, params ...interface{})
	Errorf(format string, params ...interface{})
	Panicf(format string, params ...interface{})
	DPanicf(format string, params ...interface{})
	Fatalf(format string, params ...interface{})
	Level() Level
}
```
## 核心概念
### Level 日志级别，默认级别为Info
* Debug
* Info
* Warn
* Error
* DPanic
* Panic
* Fatal
### Option 是一个配置选项，可以用来配置日志的输出格式，输出位置，日志级别等。其中nlog和xfmt的Option较为简单，xzap的Option则复杂一些。
  * xfmt的Option
	* Level 日志级别
  * nlog的Option
    * Name  日志实例名
	* Flag  对应Golang中log的Flag // log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix| log.LUTC| log.Llongfile
	* Level 日志级别
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

## 如何使用
### 使用默认配置的日志实例
 * 默认xfmt
   * 默认配置
	```golang
		DefaultName  = "Defaultxfmt"
		DefaultLevel = "INFO"
	```
	* 使用
	```golang
		import (
			"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
		)
		xfmt.Info("hello xfmt")
 * 默认nlog
   * 默认配置
	```golang
		DefaultName  = "Defaultnlog"
		DefaultFlag  = "Ldate | Ltime | Lmicroseconds | Lshortfile | Lmsgprefix"
		DefaultLevel = "INFO"
	```
	* 使用
	```golang
		import (
			"github.com/NetEase-Media/easy-ngo/xlog/nlog"
		)
		nlog.Info("hello nlog")
	```
### 使用自定义配置日志实例
第一步，添加日志的相关配置，支持yaml、toml
在 app.yaml 配置文件中添加如下配置
```yaml
ngo:
 logger:
  nlog:
   -
    Name : "nlog1"
    Level: "DEBUG"
    Flag : "Ldate | Ltime | Lmicroseconds | Lshortfile | Lmsgprefix"
  xzap: 
   -
    Name : "xzap2"
    Level: "INFO"
    NoFile : true
    Format : "text"
    WritableCaller : true
    Skip : 2
    WritableStack : false
    Path : "./logs"
    FileName : "esay-ngo"
    ErrlogLevel : "ERROR"
    ErrorPath : "./logs"
    MaxAge : 7
    MaxBackups : 7
    MaxSize : 1024
    Compress : false
```
或在 app.toml 配置文件中添加如下配置
```toml
[[ngo.logger.nlog]]
    Name = "nlog1"
    Level = "DEBUG"
    Flag = "Ldate | Ltime | Lmicroseconds | Lshortfile | Lmsgprefix"
[[ngo.logger.xzap]]
    Name = "xzap1"
    Level= "INFO"
    NoFile = true
    Format = "text"
    WritableCaller = true
    Skip = 2
    WritableStack = false
    Path = "./logs"
    FileName = "esay-ngo"
    ErrlogLevel = "ERROR"
    ErrorPath = "./logs"
    MaxAge = 7
    MaxBackups = 7
    MaxSize = 1024
    Compress = false
```
### 第二步，使用
在go的启动文件中添加如下代码
```golang
package main

import (
	"fmt"

	"github.com/NetEase-Media/easy-ngo/application"
	_ "github.com/NetEase-Media/easy-ngo/application/r/rconfig"
	"github.com/NetEase-Media/easy-ngo/application/r/rlog"
	_ "github.com/NetEase-Media/easy-ngo/application/r/rlog/rzap"
)

// go run main.go -c ./app.yaml
func main() {
	app := application.Default()
	app.Initialize()
	app.Startup()
	nlogLogger := rlog.GetLogger("nlog2")
	if nlogLogger == nil {
		fmt.Print("failed....")
		return
	}
	nlogLogger.Infof("hello world1")

	xzapLogger := rlog.GetLogger("xzap1")
	if xzapLogger == nil {
		fmt.Print("failed....")
		return
	}
	xzapLogger.Infof("hello world2")
	fmt.Print("success")
}
```




