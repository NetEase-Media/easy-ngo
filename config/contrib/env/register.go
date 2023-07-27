package env

import "github.com/NetEase-Media/easy-ngo/config"

func init() {
	config.Register(EnvConfigSourceName, New())
}
