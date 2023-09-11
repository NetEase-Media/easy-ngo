package main

import (
	"context"
	"fmt"

	"github.com/NetEase-Media/easy-ngo/app"
	xfasthttp "github.com/NetEase-Media/easy-ngo/app/plugins/plugin_xfasthttp"
)

func main() {
	app := app.New()
	app.Init()
	code, err := xfasthttp.GetXfasthttp().Get("https://www.baidu.com").Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(code)
}
