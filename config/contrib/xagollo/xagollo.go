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

package xagollo

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/philchia/agollo"
)

// this component only support toml file unmarshal
// if you want more file type to unmarshal
// implement it yourself.

const (
	scheme = "apollo"
)

type Xagollo struct {
	client      agollo.Client
	namespace   string
	propertyKey string
}

func (ag *Xagollo) Load(sourcePathes []string) (map[string]interface{}, error) {
	if len(sourcePathes) == 0 {
		return nil, nil
	}
	var all string
	for _, v := range sourcePathes {
		if len(v) == 0 {
			continue
		}
		c := parseAndCreate(v)
		if c == nil {
			continue
		}
		content := c.client.GetNameSpaceContent(c.namespace, "")
		if len(content) == 0 {
			continue
		}
		all = fmt.Sprintf("%s%s\n", all, content)
	}
	if len(all) == 0 {
		return nil, nil
	}
	M := make(map[string]interface{})
	err := toml.Unmarshal([]byte(all), &M)
	if err != nil {
		return nil, err
	}
	return M, nil
}

func parseAndCreate(sourcePathe string) *Xagollo {
	url, err := url.Parse(sourcePathe)
	if err != nil {
		return nil
	}
	if url.Scheme != scheme {
		return nil
	}

	p := url.Query()
	var nsn = strings.Split(p.Get("namespaceNames"), ",")
	clientConfig := &agollo.Conf{
		AppID:          p.Get("appId"),
		Cluster:        p.Get("cluster"),
		NameSpaceNames: nsn,
		IP:             url.Host,
	}
	client := agollo.NewClient(clientConfig)
	ag := &Xagollo{
		client:      *client,
		namespace:   nsn[0],
		propertyKey: p.Get("propertyKey"),
	}
	ag.client.Start() // 开始运行
	return ag
}
