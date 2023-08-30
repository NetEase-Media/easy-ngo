package pluginxgin

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/NetEase-Media/easy-ngo/config"
	"github.com/NetEase-Media/easy-ngo/server"
	"github.com/NetEase-Media/easy-ngo/server/contrib/xgin"
)

func init() {
	app.RegisterPlugin(app.Initialize, Initialize)
	app.RegisterPlugin(app.Starting, Serve)
	app.RegisterPlugin(app.Stopping, Shutdown)
}

func Initialize(ctx context.Context) error {
	c := server.DefaultConfig()
	if err := config.UnmarshalKey("server", c); err != nil {
		return err
	}
	WithServer(xgin.New(c))
	return GetServer().Init()
}

func Serve(ctx context.Context) error {
	return GetServer().Serve()
}

func Shutdown(ctx context.Context) error {
	return GetServer().Shutdown()
}
