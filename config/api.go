package config

import (
	"time"
)

var config Config

func GetString(key string) string {
	return config.GetString(key)
}

func GetInt(key string) int {
	return config.GetInt(key)
}

func GetBool(key string) bool {
	return config.GetBool(key)
}

func GetTime(key string) time.Time {
	return config.GetTime(key)
}

func GetFloat64(key string) float64 {
	return config.GetFloat64(key)
}

func UnmarshalKey(key string, rawVal interface{}) error {
	return config.UnmarshalKey(key, &rawVal)
}

func WithConfig(c Config) {
	config = c
}
