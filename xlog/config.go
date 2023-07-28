package xlog

type Config struct {
	Format     string //日志格式，支持json、text
	Path       string //保存日志文件的路径
	FileName   string //日志文件名
	MaxAge     int    // 保留旧文件的最大天数，默认7天
	MaxBackups int    // 保留旧文件的最大个数，默认7个
	MaxSize    int    // 在进行切割之前，日志文件的最大大小（以MB为单位）默认1024
	Compress   bool   // 是否压缩/归档旧文件
}

func DefaultConfig() *Config {
	return &Config{
		Format:     "json",
		Path:       "./logs",
		FileName:   "app.log",
		MaxAge:     7,
		MaxBackups: 7,
		MaxSize:    1024,
		Compress:   false,
	}
}
