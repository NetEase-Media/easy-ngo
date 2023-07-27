package config

import "time"

var config *Config

func GetString(key string) string {
	return config.Viper.GetString(key)
}

func GetInt(key string) int {
	return config.Viper.GetInt(key)
}

func GetBool(key string) bool {
	return config.Viper.GetBool(key)
}

func GetTime(key string) time.Time {
	return config.Viper.GetTime(key)
}

func GetFloat64(key string) float64 {
	return config.Viper.GetFloat64(key)
}

func GetDuration(key string) time.Duration {
	return config.Viper.GetDuration(key)
}

func UnmarshalKey(key string, rawVal interface{}) {
	config.Viper.UnmarshalKey(key, &rawVal)
}

func WithConfig(c *Config) {
	config = c
}
