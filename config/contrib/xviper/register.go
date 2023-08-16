package xviper

import "github.com/NetEase-Media/easy-ngo/config"

const (
	EnvConfigSourceName  = "env"
	FileConfigSourceName = "file"
)

var xviper = New()

func init() {
	config.Register(FileConfigSourceName, xviper)
	config.Register(EnvConfigSourceName, xviper)
	config.WithVendor(xviper)
}
