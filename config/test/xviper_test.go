package test

import (
	"os"
	"testing"

	"github.com/NetEase-Media/easy-ngo/config"
	_ "github.com/NetEase-Media/easy-ngo/config/contrib/xviper"
	"gotest.tools/assert"
)

type App struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
}

func TestXViper(t *testing.T) {
	os.Setenv("APP_NAME", "easy-ngo")
	os.Setenv("APP_VERSION", "v1.0.0")
	os.Setenv("APP_PORT", "8080")
	c := config.New()
	c.AddProtocol("env://prefix=APP")
	c.AddProtocol("file://type=toml;path=./;name=test2")
	c.Init()
	c.ReadConfig()
	assert.Equal(t, "v1.0.0", "v1.0.0")
	assert.Equal(t, config.GetString("name"), "easy-ngo")
	assert.Equal(t, config.GetString("version"), "v1.0.0")
	assert.Equal(t, config.GetInt("port"), 8080)
	app := &App{}
	config.UnmarshalKey("app", app)
	assert.Equal(t, app.Name, "test")
	assert.Equal(t, app.Version, "v1.0.0")
	assert.Equal(t, app.Port, 8080)
}
