package pluginxgin

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/NetEase-Media/easy-ngo/config"
	"github.com/NetEase-Media/easy-ngo/server/contrib/xgin"
)

func init() {
	app.RegisterPlugin(app.Initialize, Initialize)
	app.RegisterPlugin(app.Starting, Serve)
}

func Initialize(ctx context.Context) error {
	c := xgin.DefaultConfig()
	if err := config.UnmarshalKey("server.gin", c); err != nil {
		return err
	}
	xgin.WithServer(xgin.New(c))
	return xgin.GetServer().Init()
}

func Serve(ctx context.Context) error {
	return xgin.GetServer().Serve()
}
