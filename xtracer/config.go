package xtracer

type EXPORTER_NAME string

const (
	EXPORTER_NAME_JAEGER EXPORTER_NAME = "jaeger"
	EXPORTER_NAME_OLTP   EXPORTER_NAME = "oltp"
	EXPORTER_NAME_STDOUT EXPORTER_NAME = "stdout"
)

type Config struct {
	// 采样率
	SampleRate float64
	// 采样器
	ExporterName EXPORTER_NAME
	// OLTP采样器服务地址
	ExporterEndpoint string
	// OLTP采样器服务名称
	ServiceName string
}

func DefaultConfig() *Config {
	return &Config{
		SampleRate:   100,
		ExporterName: EXPORTER_NAME_STDOUT,
	}
}
