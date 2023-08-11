package main

import (
	"github.com/NetEase-Media/easy-ngo/server"
	"github.com/NetEase-Media/easy-ngo/server/contrib/xgin"
	"github.com/gin-gonic/gin"
)

var s server.Server

func main() {
	s = xgin.New(xgin.DefaultConfig())
	s.GET("/test", func(ctx *gin.Context) {
		ctx.String(200, "hello world")
	})
	s.Init()
	s.Serve()
}
