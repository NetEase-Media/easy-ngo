package pluginxgrom

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/NetEase-Media/easy-ngo/clients/xgorm"
	"github.com/NetEase-Media/easy-ngo/config"
)

func init() {
	app.RegisterPlugin(app.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	configs := make([]xgorm.Config, 0)
	if err := config.UnmarshalKey("gorm", configs); err != nil {
		return err
	}
	if len(configs) == 0 {
		configs = append(configs, *xgorm.DefaultConfig())
	}
	for _, config := range configs {
		cli := xgorm.New(&config)
		cli.Init()
		if err := cli.Init(); err != nil {
			return err
		}
		set(config.Name, cli)
	}
	return nil
}
