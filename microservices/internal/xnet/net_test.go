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

package xnet

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalIP(t *testing.T) {
	ip, err := LocalIP()
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)
	assert.NotEqual(t, "localhost", ip)
	assert.NotEqual(t, "127.0.0.1", ip)
}

func TestParseAddr(t *testing.T) {
	host, port := ParseAddr("localhost:8080")
	assert.Equal(t, "localhost", host)
	assert.Equal(t, "8080", port)
}

func TestParseURL(t *testing.T) {
	host, err := ParseURL("http://localhost:8080/hello")
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8080", host)
}

func TestPort(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:8080")
	assert.NoError(t, err)
	port, ok := Port(lis)
	assert.True(t, ok)
	assert.Equal(t, 8080, port)
}
