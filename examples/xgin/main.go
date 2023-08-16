package main

import (
	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/gin-gonic/gin"

	xgin "github.com/NetEase-Media/easy-ngo/app/plugins/plugin_xgin"
)

func main() {
	app := app.New()
	app.Start(addRoutes)
}

func addRoutes() error {
	xgin.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, "ddd")
	})
	return nil
}