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

package xgin

import (
	"fmt"
	"net"
	"net/http"

	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	*http.Server
	config   *Config
	listener net.Listener
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		Engine: gin.New(),
	}
}

func (server *Server) Serve() error {
	server.Server = &http.Server{
		Addr:    server.Address(),
		Handler: server,
	}
	err := server.Server.Serve(server.listener)
	if err == http.ErrServerClosed {
		xlog.Panicf("gin server closed,%v", err)
		return nil
	}
	return nil
}

func (server *Server) Shutdown() error {
	return nil
}

func (server *Server) Healthz() bool {
	return true
}

func (server *Server) Address() string {
	return fmt.Sprintf("%s:%d", server.config.Host, server.config.Port)
}
