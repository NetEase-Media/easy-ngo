package file

import (
	"strings"

	"github.com/NetEase-Media/easy-ngo/config"
)

const FileConfigSourceName = "file"

type FileConfigSource struct {
	name     string
	fileType string
	path     []string
}

func New() *FileConfigSource {
	return &FileConfigSource{
		path: make([]string, 0),
	}
}

func (e *FileConfigSource) Init(protocol string, config *config.Config) error {
	kvs := strings.Split(protocol, ";")
	for _, kv := range kvs {
		if strings.HasPrefix(kv, "name=") {
			e.name = strings.TrimPrefix(kv, "name=")
		} else if strings.HasPrefix(kv, "type=") {
			e.fileType = strings.TrimPrefix(kv, "type=")
		} else if strings.HasPrefix(kv, "path=") {
			e.path = strings.Split(strings.TrimPrefix(kv, "path="), ",")
		}
	}
	config.Viper.SetConfigName(e.name)
	config.Viper.SetConfigType(e.fileType)
	for _, p := range e.path {
		config.Viper.AddConfigPath(p)
	}
	return nil
}

func (e *FileConfigSource) Read(config *config.Config) error {
	return config.Viper.ReadInConfig()
}
