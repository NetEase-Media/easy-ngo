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
	"time"

	"github.com/NetEase-Media/easy-ngo/server"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/NetEase-Media/easy-ngo/xmetrics"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	*http.Server
	config   *server.Config
	listener net.Listener

	metrics *server.HttpMetrics
}

func New(config *server.Config) *Server {
	s := &Server{
		config: config,
		Engine: gin.New(),
	}
	if config.Metrics.Enabled {
		s.metrics = server.NewHttpMetrics(xmetrics.GetProvider(), xmetrics.Bucket(config.Metrics.Bucket))
	}
	return s
}

func (server *Server) Serve() error {
	server.Server = &http.Server{
		Addr:    server.Address(),
		Handler: server,
	}
	err := server.Server.Serve(server.listener)
	if err != nil && err == http.ErrServerClosed {
		xlog.Panicf("close gin[%s]", err)
		return nil
	}
	return nil
}

func (s *Server) Init() error {
	if s.config.Metrics.Enabled {
		s.metrics.Init()
		s.Use(s.metricsMiddleware())
	}
	//初始化Tracer
	s.initTracer()
	s.Use(s.traceMiddleware())
	listener, err := net.Listen("tcp", s.Address())
	if err != nil {
		xlog.Panicf("gin Init error![%s]", err)
		return err
	}
	s.listener = listener
	gin.SetMode(string(s.config.Mode))
	return nil
}

func (s *Server) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		s.metrics.Record(server.HttpLabels{
			Url:    c.Request.URL.Path,
			Method: c.Request.Method,
			Code:   c.Writer.Status(),
			Domain: c.Request.Host,
		}, start, time.Now())
	}
}

func (s *Server) Shutdown() error {
	return s.Server.Close()
}

func (server *Server) Healthz() bool {
	return true
}

func (server *Server) Address() string {
	return fmt.Sprintf("%s:%d", server.config.Host, server.config.Port)
}
