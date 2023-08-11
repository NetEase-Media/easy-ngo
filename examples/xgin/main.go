package main

import (
	"github.com/NetEase-Media/easy-ngo/config"
	"github.com/NetEase-Media/easy-ngo/server"
	"github.com/NetEase-Media/easy-ngo/server/contrib/xgin"
	"github.com/gin-gonic/gin"
)

var s server.Server

func main() {
	c := config.New()
	c.AddProtocol("file://type=yaml;path=./;name=app")
	c.Init()
	c.ReadConfig()
	config := xgin.DefaultConfig()
	c.UnmarshalKey("server", config)
	s = xgin.New(config)
	s.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, "ddd")
	})
	s.Init()
	s.Serve()
}
