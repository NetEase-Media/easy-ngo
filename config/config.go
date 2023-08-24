package config

import (
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	protocols map[string]string
}

func New() *Config {
	return &Config{
		protocols: make(map[string]string),
	}
}

func (c *Config) Init() error {
	//初始化Contrib
	for key, value := range c.protocols {
		//初始化Contrib
		vendors[key].Init(value)
	}
	return nil
}

func (c *Config) AddProtocol(protocol string) *Config {
	//如果没有找到://，则默认为file://，文件类型为yaml
	if !strings.Contains(protocol, "://") {
		//从后往前截取.号，.前面为文件路径和文件名，后面为文件类型
		fileNameAndExt := strings.Split(protocol, ".")
		dir, file := filepath.Split(fileNameAndExt[0])
		//从后往前截取第一个/，前面为路径，后面为文件名
		protocol = "path=" + dir + ";name=" + file
		if len(fileNameAndExt) == 2 {
			protocol = protocol + ";type=" + fileNameAndExt[1]
		}
		protocol = "file://" + protocol
	}
	c.protocols[strings.Split(protocol, "://")[0]] = protocol
	return c
}

func (c *Config) ReadConfig() error {
	for scheme, _ := range c.protocols {
		if "file" == scheme {
			err := vendors[scheme].Read()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Config) GetString(key string) string {
	return propertiesVendor.GetString(key)
}

func (c *Config) UnmarshalKey(key string, rawVal interface{}) error {
	return propertiesVendor.UnmarshalKey(key, rawVal)
}

func (c *Config) GetInt(key string) int {
	return propertiesVendor.GetInt(key)
}

func (c *Config) GetBool(key string) bool {
	return propertiesVendor.GetBool(key)
}

func (c *Config) GetTime(key string) time.Time {
	return propertiesVendor.GetTime(key)
}

func (c *Config) GetFloat64(key string) float64 {
	return propertiesVendor.GetFloat64(key)
}
