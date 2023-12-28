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
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	EnvConfigSourceName  = "env"
	FileConfigSourceName = "file"
)

type XViper struct {
	viper *viper.Viper
}

func New() Config {
	return &XViper{
		viper: viper.New(),
	}
}

func (xviper *XViper) Init(protocols ...string) error {
	for _, protocol := range protocols {
		scheme := protocol[:strings.Index(protocol, "://")]
		switch scheme {
		case EnvConfigSourceName:
			xviper.initEnv(protocol[strings.Index(protocol, "://")+3:])
		case FileConfigSourceName:
			if err := xviper.initFile(protocol[strings.Index(protocol, "://")+3:]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (xviper *XViper) initEnv(protocol string) {
	kvs := strings.Split(protocol, ";")
	for _, kv := range kvs {
		if strings.HasPrefix(kv, "prefix=") {
			xviper.viper.SetEnvPrefix(strings.TrimPrefix(kv, "prefix="))
		}
	}
	xviper.viper.AutomaticEnv()
}

func (xviper *XViper) initFile(protocol string) error {
	kvs := strings.Split(protocol, ";")
	for _, kv := range kvs {
		if strings.HasPrefix(kv, "name=") {
			xviper.viper.SetConfigName(strings.TrimPrefix(kv, "name="))
		} else if strings.HasPrefix(kv, "type=") {
			xviper.viper.SetConfigType(strings.TrimPrefix(kv, "type="))
		} else if strings.HasPrefix(kv, "path=") {
			paths := strings.Split(strings.TrimPrefix(kv, "path="), ",")
			for _, path := range paths {
				xviper.viper.AddConfigPath(path)
			}
		}
	}
	return xviper.viper.ReadInConfig()
}

func (xviper *XViper) Get(key string) interface{} {
	return xviper.viper.Get(key)
}

func (xviper *XViper) GetString(key string) string {
	return xviper.viper.GetString(key)
}

func (xviper *XViper) UnmarshalKey(key string, rawVal interface{}) error {
	return xviper.viper.UnmarshalKey(key, rawVal)
}

func (xviper *XViper) GetInt(key string) int {
	return xviper.viper.GetInt(key)
}

func (xviper *XViper) GetBool(key string) bool {
	return xviper.viper.GetBool(key)
}

func (xviper *XViper) GetTime(key string) time.Time {
	return xviper.viper.GetTime(key)
}

func (xviper *XViper) GetFloat64(key string) float64 {
	return xviper.viper.GetFloat64(key)
}
