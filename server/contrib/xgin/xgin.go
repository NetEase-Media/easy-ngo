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
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	*http.Server
	config   *Config
	listener net.Listener

	httpMetrics *server.HttpMetrics
	routes      []*server.Route
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		Engine: gin.New(),
		routes: make([]*server.Route, 0),
	}
}

func (server *Server) Serve() error {
	server.Server = &http.Server{
		Addr:    server.Address(),
		Handler: server,
	}
	return server.Server.Serve(server.listener)
}

func (s *Server) Init() error {
	if s.config.EnabledMetric {
		s.httpMetrics = server.NewHttpMetrics()
		s.httpMetrics.Init()
		s.Use(s.metricsMiddleware())
	}
	if s.config.EnabledTrace {
		s.Use(s.traceMiddleware())
	}
	for _, route := range s.routes {
		switch route.Method {
		case server.GET:
			s.Engine.GET(route.RelativePath, s.handlerWrapper(route.Handler))
		case server.POST:
			s.Engine.POST(route.RelativePath, s.handlerWrapper(route.Handler))
		case server.PUT:
			s.Engine.PUT(route.RelativePath, s.handlerWrapper(route.Handler))
		case server.DELETE:
			s.Engine.DELETE(route.RelativePath, s.handlerWrapper(route.Handler))
		case server.PATCH:
			s.Engine.PATCH(route.RelativePath, s.handlerWrapper(route.Handler))
		case server.HEAD:
			s.Engine.HEAD(route.RelativePath, s.handlerWrapper(route.Handler))
		case server.OPTIONS:
			s.Engine.OPTIONS(route.RelativePath, s.handlerWrapper(route.Handler))
		}
	}
	listener, err := net.Listen("tcp", s.Address())
	if err != nil {
		return err
	}
	s.listener = listener
	gin.SetMode(string(s.config.Mode))
	return nil
}

func (s *Server) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if s.config.MetricsPath == c.Request.URL.Path {
		// 	c.Next()
		// 	return
		// }
		start := time.Now()
		c.Next()
		s.httpMetrics.Record((time.Now().Nanosecond()-start.Nanosecond())/1e6, server.HttpLabels{
			Url:    c.Request.URL.Path,
			Method: c.Request.Method,
			Code:   c.Writer.Status(),
			Domain: c.Request.Host,
		})
	}
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

func (s *Server) PUT(relativePath string, handler any) {
	s.appendRoute(server.PUT, relativePath, handler)
}

func (s *Server) GET(relativePath string, handler any) {
	s.appendRoute(server.GET, relativePath, handler)
}

func (s *Server) POST(relativePath string, handler any) {
	s.appendRoute(server.POST, relativePath, handler)
}

func (s *Server) DELETE(relativePath string, handler any) {
	s.appendRoute(server.POST, relativePath, handler)
}

func (s *Server) PATCH(relativePath string, handler any) {
	s.appendRoute(server.POST, relativePath, handler)
}

func (s *Server) HEAD(relativePath string, handler any) {
	s.appendRoute(server.POST, relativePath, handler)
}

func (s *Server) OPTIONS(relativePath string, handler any) {
	s.appendRoute(server.POST, relativePath, handler)
}

func (s *Server) appendRoute(method server.METHOD, relativePath string, handler any) {
	s.routes = append(s.routes, &server.Route{
		Method:       method,
		RelativePath: relativePath,
		Handler:      handler,
	})
}

func (s *Server) handlerWrapper(handler any) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.(func(c *gin.Context))(c)
	}
}
