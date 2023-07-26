package xecho

import (
	"net"
	"strconv"

	"github.com/NetEase-Media/easy-ngo/server"
	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo

	config *Config
}

func New(config *Config) *Server {
	return &Server{
		Echo: echo.New(),
	}
}

func (s *Server) Serve() error {
	err := s.Echo.Start(net.JoinHostPort(s.config.Host, strconv.Itoa(s.config.Port)))
	return err
}

func (s *Server) Shutdown() error {
	return nil
}

func (*Server) Healthz() bool {
	return true
}

func (*Server) Info() *server.Info {
	return &server.Info{}
}
