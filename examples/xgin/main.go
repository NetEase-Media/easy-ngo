package main

import (
	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/gin-gonic/gin"

	_ "github.com/NetEase-Media/easy-ngo/app/plugins/plugin_xgin"
	"github.com/NetEase-Media/easy-ngo/server/contrib/xgin"
)

func main() {
	app := app.New()
	err := app.Init(addRoutes)
	if err != nil {
		panic(err)
	}
	app.Start()
}

func addRoutes() error {
	xgin.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, "ddd")
	})
	return nil
}
