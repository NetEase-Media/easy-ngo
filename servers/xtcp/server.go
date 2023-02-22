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

	"github.com/NetEase-Media/easy-ngo/xlog"
)

type Server struct {
	listener net.Listener
	opt      *Option
	logger   xlog.Logger
	Handler  func(net.Conn, context.Context)
}

func New(opt *Option) *Server {
	if opt == nil {
		opt = defaultOption()
	}
	return &Server{
		opt: opt,
	}
}

func defaultOption() *Option {
	return &Option{
		IP:   "0.0.0.0",
		Port: 8888,
	}
}

func (server *Server) Initial() (err error) {
	add := fmt.Sprintf("%s:%d", server.opt.IP, server.opt.Port)
	listen, err := net.Listen("tcp", add)
	if err != nil {
		return
	}
	server.listener = listen
	return nil
}

func (server *Server) RegisterHandler(handler func(net.Conn, context.Context)) {
	server.Handler = handler
}

func (server *Server) Listen() error {
	if server.Handler == nil {
		if server.logger != nil {
			server.logger.Errorf("server listen failed for no handler.")
		}
		return NoHandlerError // no handler error
	}
	defer server.listener.Close()
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			if server.logger != nil {
				server.logger.Errorf("tcp listener accept err, %s", err.Error())
			}
			continue // always listening
		}
		if server.logger.Level() == xlog.DebugLevel && server.logger != nil {
			server.logger.Debugf("accept connetion, remote address %v", conn.RemoteAddr())
		}
		go server.Handler(conn, context.Background())
	}
}
