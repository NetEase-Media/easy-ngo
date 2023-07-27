package test

import (
	"os"
	"testing"

	"github.com/NetEase-Media/easy-ngo/config"
	_ "github.com/NetEase-Media/easy-ngo/config/contrib/env"
	"github.com/go-playground/assert/v2"
)

func TestEnv(t *testing.T) {
	os.Setenv("APP_NAME", "easy-ngo")
	os.Setenv("APP_VERSION", "v1.0.0")
	c := config.New()
	c.AddProtocol("env://prefix=APP")
	c.Init()
	config.WithConfig(c)
	assert.Equal(t, config.GetString("name"), "easy-ngo")
	assert.Equal(t, config.GetString("version"), "v1.0.0")
}
