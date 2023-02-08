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

package xviper

import (
	"github.com/spf13/viper"
)

type Xviper struct {
	*viper.Viper
}

func New() *Xviper {
	return &Xviper{
		Viper: viper.New(),
	}
}

func (vip *Xviper) Load(sourcePathes []string) (map[string]interface{}, error) {
	if len(sourcePathes) == 0 {
		return nil, nil
	}
	for _, v := range sourcePathes {
		tvip := &Xviper{
			Viper: viper.New(),
		}
		m, err := tvip.load(v)
		if err != nil {
			return nil, err
		}
		vip.MergeConfigMap(m)
	}
	M := vip.AllSettings()
	return M, nil
}

func (vip *Xviper) load(sourcePath string) (map[string]interface{}, error) {
	if len(sourcePath) == 0 {
		return nil, nil
	}
	vip.SetConfigFile(sourcePath)
	vip.ReadInConfig()
	// M := make(map[string]interface{})
	M := vip.AllSettings()
	return M, nil
}
