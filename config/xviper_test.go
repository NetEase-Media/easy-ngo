// Copyright 2022 NetEase Media Technology（Beijing）Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"os"
	"testing"

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
	c := New()
	c.Init("env://prefix=APP", "file://type=toml;path=./file;name=test2")
	WithConfig(c)
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
