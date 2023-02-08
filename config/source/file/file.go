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

package file

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

const (
	defaultConfigPath = "./application.toml"
)

// local or remote ?
type File struct {
	Path string
	// Unmarshaller Unmarshaller
	EnableWatch bool
}

// if we want abstract the file ...
// we need it
type Unmarshaller = func([]byte, interface{}) error

// we just implement local file read and format so far.
func (file *File) Read() (content []byte, err error) {
	content, err = ioutil.ReadFile(file.Path)
	return
}

// build-in file path, can not be modified.
type defaultConfigFile struct {
	*TomlFile
}

// for framework to use.
// we may want some default configuration.
func Load(sourcePathes []string) (map[string]interface{}, error) {
	f := &defaultConfigFile{
		TomlFile: &TomlFile{
			File: &File{
				Path: defaultConfigPath,
			},
		},
	}
	_, err := os.Stat(defaultConfigPath)
	if err != nil && os.IsNotExist(err) {
		return nil, nil
	}
	return f.Load(sourcePathes)
}

type LocalFile struct {
	*File
}

func (lf *LocalFile) Load(sourcePathes []string) (map[string]interface{}, error) {
	var pathes []string
	for _, v := range sourcePathes {
		if len(v) == 0 {
			continue
		}
		urlObj, err := url.Parse(v)
		if err != nil {
			continue
		}
		if len(urlObj.Scheme) == 0 { // local file
			pathes = append(pathes, urlObj.Path)
		}
	}
	if len(pathes) == 0 {
		return nil, nil
	}
	M := make(map[string]interface{})
	var err error
	for _, v := range pathes {
		if len(v) == 0 {
			continue
		}
		var m map[string]interface{}
		if strings.HasSuffix(v, ".toml") {
			parser := TomlFile{
				File: &File{
					Path: v,
				},
			}
			m, err = parser.Load(nil)
			if err != nil {
				fmt.Printf("toml load file:%s error:%s\n", v, err)
				continue
			}
		} else if strings.HasSuffix(v, ".yaml") {
			parser := YamlFile{
				File: &File{
					Path: v,
				},
			}
			m, err = parser.Load(nil)
			if err != nil {
				fmt.Printf("yaml load file:%s error:%s\n", v, err)
				continue
			}
		} else if strings.HasSuffix(v, ".json") {
			parser := JsonFile{
				File: &File{
					Path: v,
				},
			}
			m, err = parser.Load(nil)
			if err != nil {
				fmt.Printf("json load file:%s error:%s\n", v, err)
				continue
			}
		}
		if len(m) == 0 {
			continue
		}
		for k := range m {
			M[k] = m[k] //
		}
	}
	return M, nil
}
