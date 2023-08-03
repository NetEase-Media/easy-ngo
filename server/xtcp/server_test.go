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

package xtcp

import (
	"context"
	"fmt"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	opt := Option{
		Name: "server01",
		IP:   "0.0.0.0",
		Port: 8888,
	}

	server := New(&opt)
	err := server.Initial()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	server.RegisterHandler(func(con net.Conn, ctx context.Context) {
		var buf []byte = make([]byte, 1024)
		for {
			_, err := con.Read(buf)
			if err != nil {
				return
			}
			fmt.Print(string(buf))
		}
	})
	server.Listen()
}
