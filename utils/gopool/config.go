package gopool

const (
	defaultScalaThreshold = 1
)

// Config is used to config pool.
type Config struct {
	// threshold for scale.
	// new goroutine is created if len(task chan) > ScaleThreshold.
	// defaults to defaultScalaThreshold.
	ScaleThreshold int32
}

// NewConfig creates a default Config.
func NewConfig() *Config {
	c := &Config{
		ScaleThreshold: defaultScalaThreshold,
	}
	return c
}
