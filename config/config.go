package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/NetEase-Media/easy-ngo/config/source"
	"github.com/NetEase-Media/easy-ngog/source/env"
	"github.com/NetEase-Media/easy-ngog/source/file"
	"github.com/NetEase-Media/easy-ngog/source/parameter"
)

type Config struct {
	KeyMap       *sync.Map
	sources      []source.ConfigSource
	configPathes []string
}

func New(configSource string) *Config {
	var sourcePathes []string
	if len(configSource) == 0 { // read from default file
		sourcePathes = loadDefaultSourcePath()
	} else { // read from -c parameter
		sourcePathes = strings.Split(configSource, ";")
	}
	// if len(sourcePathes) == 0 {
	// 	panic("App need default config source. ^@^")
	// }
	// default implement
	environments := env.New()
	parameters := parameter.New()
	localFile := &file.LocalFile{}
	c := &Config{KeyMap: &sync.Map{}, sources: []source.ConfigSource{environments, parameters, localFile}}
	c.configPathes = sourcePathes
	return c
}

func NewDefault() *Config {
	// default implement
	environments := env.New()
	parameters := parameter.New()
	localFile := &file.LocalFile{}
	c := &Config{KeyMap: &sync.Map{}, sources: []source.ConfigSource{environments, parameters, localFile}}
	return c
}

func (c *Config) ParseAndSetSourePath(configSource string) {
	var sourcePathes []string
	if len(configSource) == 0 { // read from default file
		sourcePathes = loadDefaultSourcePath()
	} else { // read from -c parameter
		sourcePathes = strings.Split(configSource, ";")
	}
	c.configPathes = sourcePathes
}

func loadDefaultSourcePath() []string {
	m, err := file.Load(nil)
	var path interface{}
	if err == nil && m != nil {
		M := Flattening(m)
		path = M["default.config.addr"]
	}
	if path == nil {
		return nil
	}
	switch path.(type) {
	case string:
		return strings.Split(path.(string), ";")
	case []interface{}:
		var ret []string
		for _, v := range path.([]interface{}) {
			ret = append(ret, v.(string))
		}
		return ret
	default:
		return nil
	}
}

func (c *Config) Initialize() error {
	// load default application config
	if len(c.configPathes) == 0 {
		m, err := file.Load(nil)
		if err != nil {
			return err
		}
		M := Flattening(m)
		if len(M) == 0 {
			return nil
		}
		for k, v := range M {
			c.KeyMap.Store(k, v)
		}
		return nil
	}
	return c.load()
}

func (c *Config) register(src source.ConfigSource) {
	c.sources = append(c.sources, src)
}

func (c *Config) load() error {
	M := make(map[string]interface{})
	for _, source := range c.sources {
		if source == nil {
			continue
		}
		m, err := source.Load(c.configPathes)
		if err != nil {
			return err
		}
		if m == nil {
			continue
		}
		expandM := Expand(m)
		bs, _ := json.MarshalIndent(expandM, "", "  ")
		fmt.Printf("source type:%s, content:%s\n", reflect.TypeOf(source), bs)
		flattenM := Flattening(expandM)
		Merge(M, flattenM)
	}
	if len(M) > 0 {
		for k, v := range M {
			c.KeyMap.Store(k, v)
		}
	}
	return nil
}
