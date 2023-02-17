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
	"net"
	"net/http"

	"github.com/NetEase-Media/easy-ngo/observability/metrics"
	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
	"github.com/NetEase-Media/easy-ngo/xlog"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	*http.Server
	option   *Option
	listener net.Listener

	Logger  xlog.Logger
	Metrics metrics.Provider
	Tracer  tracer.Provider

	// JwtAuth *JwtAuth
}

func New(option *Option, logger xlog.Logger, metrics metrics.Provider, tracer tracer.Provider) *Server {
	server := &Server{
		option:  option,
		Logger:  logger,
		Metrics: metrics,
		Tracer:  tracer,
	}
	server.Engine = gin.New()
	return server
}

func Default() *Server {
	return New(DefaultOption(), nil, nil, nil)
}

func (server *Server) Serve() error {
	server.Server = &http.Server{
		Addr:    server.option.Address(),
		Handler: server,
	}
	var err error
	if server.option.TLS != nil && server.option.TLS.Enable {
		err = server.Server.ServeTLS(server.listener, server.option.TLS.Cert, server.option.TLS.Key)
	} else {
		err = server.Server.Serve(server.listener)
	}
	if err == http.ErrServerClosed {
		server.Logger.Panicf("close gin[%s]", err)
		return nil
	}
	return nil
}

func (server *Server) Shutdown() error {
	return nil
}

func (server *Server) Initialize() error {
	if server.option.EnabledMetric {
		server.initMetrics()
		server.Use(server.metricsMiddleware())
	}
	if server.option.EnableTracer {
		server.initTracer()
		server.Use(server.traceMiddleware())
	}

	// var jwt *JwtAuth
	// // 后台登录mw
	// if server.option.Middlewares.JwtAuth.Enabled {
	// 	jwt = JwtAuthInit(server.option.Middlewares.JwtAuth)
	// 	server.Engine.Use(jwt.MiddlewareFunc())
	// 	auth := server.Engine.Group(server.option.Middlewares.JwtAuth.RoutePathPrefix + "/auth")
	// 	auth.GET("/access-token", jwt.CreateTokenHandler)
	// 	auth.GET("/refresh-token", jwt.RefreshTokenHandler)
	// }
	// server.JwtAuth = jwt

	server.Engine.Use(OutermostRecover())
	if server.option.Middlewares.AccessLog.Enabled {
		server.Engine.Use(AccessLogMiddleware(server.option.Middlewares.AccessLog))
	}
	server.Engine.Use(ServerRecover())
	server.Engine.Use(SemicolonMiddleware())
	listener, err := net.Listen("tcp", server.option.Address())
	if err != nil {
		server.Logger.Panicf("new exgin server err[%s]", err.Error())
		return err
	}
	gin.SetMode(server.option.Mode)
	server.listener = listener
	return nil
}

// AddAdminAuthHandler 扩展处理方法
// param auth openid通过后的回调，用来验证业务后台用户逻辑，返回需要存入token的信息，返回的信息尽量少，否则token很长
// param gtRsp 获取token和刷新token后回调的方法，用来自定义响应报文格式
// param unauthRsp  token验证失败的回调方法，用来自定义响应报文格式
// func (s *Server) AddAdminAuthHandler(auth Authenticator, gtRsp GenTokenResponse, unauthRsp UnauthenticatedResponse) *Server {
// 	if !s.option.Middlewares.JwtAuth.Enabled {
// 		s.Logger.Warnf("jwt auth is disabled")
// 		return s
// 	}
// 	if auth != nil {
// 		s.JwtAuth.Authenticator = auth
// 	}
// 	if gtRsp != nil {
// 		s.JwtAuth.GenTokenResponse = gtRsp
// 	}
// 	if unauthRsp != nil {
// 		s.JwtAuth.UnauthenticatedResponse = unauthRsp
// 	}
// 	return s
// }

func (s *Server) AddRoute(method Method, path string, handlers ...gin.HandlerFunc) *Server {
	s.Handle(string(method), path, handlers...)
	return s
}

func (s *Server) AddRouteWithMethods(methods []Method, path string, handlers ...gin.HandlerFunc) *Server {
	if len(methods) == 0 {
		s.Logger.Panicf("methods can not be empty")
	}
	for i := range methods {
		s.Handle(string(methods[i]), path, handlers...)
	}
	return s
}
