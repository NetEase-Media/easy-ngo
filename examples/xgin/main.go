package main

import (
	"github.com/NetEase-Media/easy-ngo/app"

	_ "github.com/NetEase-Media/easy-ngo/config/contrib/xviper"
)

func main() {
	app := app.New()
	app.Init()
	// c := config.New()
	// c.AddProtocol("file://type=yaml;path=./;name=app")
	// c.Init()
	// c.ReadConfig()
	// config := xgin.DefaultConfig()
	// c.UnmarshalKey("server", config)
	// s = xgin.New(config)
	// s.GET("/test", func(ctx *gin.Context) {
	// 	ctx.JSON(200, "ddd")
	// })
	// s.Init()
	// s.Serve()
}
