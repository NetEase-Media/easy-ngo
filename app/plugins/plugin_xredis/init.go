package pluginxredis

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/NetEase-Media/easy-ngo/clients/xredis"
	"github.com/NetEase-Media/easy-ngo/config"
)

func init() {
	app.RegisterPlugin(app.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	configs := make([]xredis.Config, 0)
	if err := config.UnmarshalKey("redis", configs); err != nil {
		return err
	}
	if len(configs) == 0 {
		configs = append(configs, *xredis.DefaultConfig())
	}
	for _, config := range configs {
		cli, _ := xredis.New(&config)
		set(config.Name, cli)
	}
	return nil
}
