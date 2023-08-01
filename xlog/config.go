package xlog

type FORMAT string
type LEVEL string

const (
	JSON FORMAT = "json"
	TEXT        = "text"
)

const (
	DEBUG LEVEL = "debug"
	INFO        = "info"
	WARN        = "warn"
	ERROR       = "error"
	PANIC       = "panic"
	FATAL       = "fatal"
)

const (
	TIMEFORMAT = "2006-01-02 15:04:05.999"
)

type Config struct {
	Format     FORMAT //日志格式，支持json、text
	Level      LEVEL  //要记录的日志级别，支持debug、info、warn、error、panic、fatal
	Path       string //保存日志文件的路径
	ErrorPath  string // 错误日志文件路径
	FileName   string //日志文件名
	MaxAge     int    // 保留旧文件的最大天数，默认7天
	MaxBackups int    // 保留旧文件的最大个数，默认7个
	MaxSize    int    // 在进行切割之前，日志文件的最大大小（以MB为单位）默认1024
	Compress   bool   // 是否压缩/归档旧文件
	Pattern    string // 日志格式的正则表达式
}

func DefaultConfig() *Config {
	return &Config{
		Format:     JSON,
		Level:      INFO,
		Path:       "./logs",
		FileName:   "app.log",
		MaxAge:     7,
		MaxBackups: 7,
		MaxSize:    1024,
		Compress:   false,
	}
}
