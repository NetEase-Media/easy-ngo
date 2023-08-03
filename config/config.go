package config

import (
	"strings"
	"time"
)

var (
	vendor Properties
)

type Config struct {
	//key: scheme, value: url
	protocols map[string]string
}

func New() *Config {
	config = &Config{
		protocols: make(map[string]string),
	}
	return config
}

func (c *Config) Init() error {
	//初始化Contrib
	for key, value := range c.protocols {
		//初始化Contrib
		configsourceBuilders[key].Init(value)
	}
	return nil
}

func (c *Config) AddProtocol(protocol string) *Config {
	c.protocols[strings.Split(protocol, "://")[0]] = protocol
	return c
}

func (c *Config) ReadConfig() error {
	for scheme, _ := range c.protocols {
		err := configsourceBuilders[scheme].Read()
		if err != nil {
			return err
		}
	}
	return nil
}

func WithVendor(v Properties) {
	vendor = v
}

func (c *Config) GetString(key string) string {
	return vendor.GetString(key)
}

func (c *Config) UnmarshalKey(key string, rawVal interface{}) error {
	return vendor.UnmarshalKey(key, rawVal)
}

func (c *Config) GetInt(key string) int {
	return vendor.GetInt(key)
}

func (c *Config) GetBool(key string) bool {
	return vendor.GetBool(key)
}

func (c *Config) GetTime(key string) time.Time {
	return vendor.GetTime(key)
}

func (c *Config) GetFloat64(key string) float64 {
	return vendor.GetFloat64(key)
}
