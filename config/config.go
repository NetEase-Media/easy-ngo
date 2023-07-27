package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	//key: scheme, value: url
	protocols map[string]string

	Viper *viper.Viper
}

func New() *Config {
	return &Config{
		protocols: make(map[string]string),
		Viper:     viper.New(),
	}
}

func (c *Config) Init() error {
	//初始化Contrib
	for key, value := range c.protocols {
		//初始化Contrib
		configsourceBuilders[key].Init(value, c)
	}
	return nil
}

func (c *Config) ParseArgument() *Config {
	return c
}

func (c *Config) AddProtocol(protocol string) *Config {
	c.protocols[strings.Split(protocol, "://")[0]] = strings.Split(protocol, "://")[1]
	return c
}

func (c *Config) ReadConfig() error {
	for scheme, _ := range c.protocols {
		err := configsourceBuilders[scheme].Read(c)
		if err != nil {
			return err
		}
	}
	return nil
}
