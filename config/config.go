package config

import "time"

type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetTime(key string) time.Time
	GetFloat64(key string) float64
	UnmarshalKey(key string, rawVal interface{}) error

	Init(protocols ...string) error
}
