package xgin

import (
	"github.com/gin-gonic/gin"
)

var s *Server

func WithServer(s1 *Server) {
	s = s1
}

func GetServer() *Server {
	return s
}

func PUT(relativePath string, handler gin.HandlerFunc) {
	s.Engine.PUT(relativePath, handler)
}

func GET(relativePath string, handler gin.HandlerFunc) error {
	s.Engine.GET(relativePath, handler)
	return nil
}

func POST(relativePath string, handler gin.HandlerFunc) {
	s.Engine.POST(relativePath, handler)
}

func DELETE(relativePath string, handler gin.HandlerFunc) {
	s.Engine.DELETE(relativePath, handler)
}

func PATCH(relativePath string, handler gin.HandlerFunc) {
	s.Engine.PATCH(relativePath, handler)
}

func HEAD(relativePath string, handler gin.HandlerFunc) {
	s.Engine.HEAD(relativePath, handler)
}

func OPTIONS(relativePath string, handler gin.HandlerFunc) {
	s.Engine.OPTIONS(relativePath, handler)
}
