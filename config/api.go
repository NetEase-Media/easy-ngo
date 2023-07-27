package config

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

func WithConfig(c *Config) {
	config = c
}
