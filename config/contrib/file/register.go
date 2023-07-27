package file

import "github.com/NetEase-Media/easy-ngo/config"

func init() {
	config.Register(FileConfigSourceName, New())
}
