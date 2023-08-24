package pluginxzk

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/app"
	config "github.com/NetEase-Media/easy-ngo/app/plugins/plugin_config"
	"github.com/NetEase-Media/easy-ngo/clients/xzk"
)

func init() {
	app.RegisterPlugin(app.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	configs := make([]xzk.Config, 0)
	if err := config.UnmarshalKey("zk", configs); err != nil {
		return err
	}
	if len(configs) == 0 {
		configs = append(configs, *xzk.DefaultConfig())
	}
	for _, config := range configs {
		cli, err := xzk.New(&config)
		if err != nil {
			return err
		}
		set(config.Name, cli)
	}
	return nil
}
