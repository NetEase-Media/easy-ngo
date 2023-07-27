package test

import (
	"testing"

	"github.com/NetEase-Media/easy-ngo/config"
	"github.com/go-playground/assert/v2"

	_ "github.com/NetEase-Media/easy-ngo/config/contrib/file"
)

type App struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
}

func TestYaml(t *testing.T) {
	c := config.New()
	c.AddProtocol("file://type=yaml;path=./;name=test1")
	c.Init()
	c.ReadConfig()
	config.WithConfig(c)
	assert.Equal(t, config.GetString("app.name"), "test")
	assert.Equal(t, config.GetString("app.version"), "v1.0.0")
	assert.Equal(t, config.GetInt("app.port"), 8080)
}

func TestToml(t *testing.T) {
	c := config.New()
	c.AddProtocol("file://type=toml;path=./;name=test2")
	c.Init()
	c.ReadConfig()
	config.WithConfig(c)
	assert.Equal(t, config.GetString("app.name"), "test")
	assert.Equal(t, config.GetString("app.version"), "v1.0.0")
	assert.Equal(t, config.GetInt("app.port"), 8080)
}

func TestStruct(t *testing.T) {
	c := config.New()
	c.AddProtocol("file://type=toml;path=./;name=test2")
	c.Init()
	c.ReadConfig()
	config.WithConfig(c)
	app := &App{}
	config.UnmarshalKey("app", app)
	assert.Equal(t, app.Name, "test")
	assert.Equal(t, app.Version, "v1.0.0")
	assert.Equal(t, app.Port, 8080)
}
