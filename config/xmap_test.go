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
	"encoding/json"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestFlatteningJson(t *testing.T) {
	data := `
	{
		"a":1,
		"b":{
			"c":2,
			"d":3,
			"e":{
				"f":4
			}
		}
	}
	`
	m := make(map[string]interface{})
	json.Unmarshal([]byte(data), &m)
	result := Flattening(m)
	t.Log(result["a"])
}

func TestFlatteningToml(t *testing.T) {
	data := `
		# Some comments.
		[alpha]
			ip = "10.0.0.1"
		[alpha.config]
			Ports = [ 8001, 8002 ]
			Location = "Toronto"
			Created = 1987-07-05T05:45:00Z
		[beta]
			ip = "10.0.0.2"
			[beta.config]
			Ports = [ 9001, 9002 ]
			Location = "New Jersey"
			Created = 1887-01-05T05:55:00Z
	`
	m := make(map[string]interface{})
	toml.Unmarshal([]byte(data), &m)
	result := Flattening(m)
	t.Log(result["a"])
}

func TestFlatteningYaml(t *testing.T) {
	var data = `
    a: Easy!
    b:
        c: 2
        d: [3, 4]
`
	m := make(map[string]interface{})
	yaml.Unmarshal([]byte(data), &m)
	result := Flattening(m)
	t.Log(result["a"])
}

func TestExpandWithoutExpr(t *testing.T) {
	data := `
		# Some comments.
		[alpha]
			ip = "10.0.0.1"
		[alpha.config]
			Ports = [ 8001, 8002 ]
			Location = "Toronto"
			Created = 1987-07-05T05:45:00Z
		[beta]
			ip = "10.0.0.2"
			[beta.config]
			Ports = [ 9001, 9002 ]
			Location = "New Jersey"
			Created = 1887-01-05T05:55:00Z
	`
	m := make(map[string]interface{})
	toml.Unmarshal([]byte(data), &m)
	result := Expand(m)
	assert.Equal(t, result, m)
}

func TestExpandStringExpr(t *testing.T) {
	data := `
		# Some comments.
		[alpha]
			ip = "10.0.0.1"
		[alpha.config]
			Ports = [ 8001, 8002 ]
			Location = "Toronto"
			Created = 1987-07-05T05:45:00Z
		[beta]
			ip = "${alpha.ip}"
			[beta.config]
			Ports = [ 9001, 9002 ]
			Location = "New Jersey"
			Created = 1887-01-05T05:55:00Z
	`
	m := make(map[string]interface{})
	toml.Unmarshal([]byte(data), &m)
	result := Expand(m)
	t.Log(m)
	t.Log(result)
	alphaIp := result["alpha"].(map[string]interface{})["ip"]
	betaIp := result["beta"].(map[string]interface{})["ip"]
	assert.Equal(t, alphaIp, betaIp)
}

func TestExpandWrongStringExpr(t *testing.T) {
	data := `
		# Some comments.
		[alpha]
			ip = "10.0.0.1"
		[alpha.config]
			Ports = [ 8001, 8002 ]
			Location = "Toronto"
			Created = 1987-07-05T05:45:00Z
		[beta]
			ip = "${alpha.ip1}"
			[beta.config]
			Ports = [ 9001, 9002 ]
			Location = "New Jersey"
			Created = 1887-01-05T05:55:00Z
	`
	m := make(map[string]interface{})
	toml.Unmarshal([]byte(data), &m)
	result := Expand(m)
	t.Log(m)
	t.Log(result)
	alphaIp := result["alpha"].(map[string]interface{})["ip"]
	betaIp := result["beta"].(map[string]interface{})["ip"]
	assert.NotEqual(t, alphaIp, betaIp)
}

func TestExpandNumberInStringExpr(t *testing.T) {
	data := `
		# Some comments.
		[alpha]
			ip = "10.0.0.1"
			port = 8080
		[alpha.config]
			Ports = [ 8001, 8002 ]
			Location = "Toronto"
			Created = 1987-07-05T05:45:00Z
		[beta]
			ip = "${alpha.port}"
			[beta.config]
			Ports = [ 9001, 9002 ]
			Location = "New Jersey"
			Created = 1887-01-05T05:55:00Z
	`
	m := make(map[string]interface{})
	toml.Unmarshal([]byte(data), &m)
	result := Expand(m)
	t.Log(m)
	t.Log(result)
	alphaIp := result["alpha"].(map[string]interface{})["ip"]
	betaIp := result["beta"].(map[string]interface{})["ip"]
	assert.NotEqual(t, alphaIp, betaIp)
}

func TestExpandDateTimeInStringExpr(t *testing.T) {
	data := `
		# Some comments.
		[alpha]
			ip = "10.0.0.1"
			port = 8080
		[alpha.config]
			Ports = [ 8001, 8002 ]
			Location = "Toronto"
			Created = 1987-07-05T05:45:00Z
		[beta]
			ip = "${alpha.config.Created}"
			[beta.config]
			Ports = [ 9001, 9002 ]
			Location = "New Jersey"
			Created = 1887-01-05T05:55:00Z
	`
	m := make(map[string]interface{})
	toml.Unmarshal([]byte(data), &m)
	result := Expand(m)
	t.Log(m)
	t.Log(result)
	alphaIp := result["alpha"].(map[string]interface{})["ip"]
	betaIp := result["beta"].(map[string]interface{})["ip"]
	assert.NotEqual(t, alphaIp, betaIp)
}

func TestExpandStringWithoutQuoteInStringExpr(t *testing.T) {
	data := `
		# Some comments.
		[alpha]
			ip = "10.0.0.1"
			port = 8080
		[alpha.config]
			Ports = [ 8001, 8002 ]
			Location = "Toronto"
			Created = 1987-07-05T05:45:00Z
		[beta]
			ip = '${alpha.config.Created}'
			[beta.config]
			Ports = [ 9001, 9002 ]
			Location = "New Jersey"
			Created = 1887-01-05T05:55:00Z
	`
	m := make(map[string]interface{})
	err := toml.Unmarshal([]byte(data), &m)
	assert.Equal(t, nil, err)
	t.Log("After toml,val", m)
	result := Expand(m)
	t.Log("After expand,val", m)
	alphaIp := result["alpha"].(map[string]interface{})["ip"]
	betaIp := result["beta"].(map[string]interface{})["ip"]
	assert.NotEqual(t, alphaIp, betaIp)
}
