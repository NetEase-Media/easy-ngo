package pluginxxxljob

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/NetEase-Media/easy-ngo/clients/xxxljob"
	"github.com/NetEase-Media/easy-ngo/config"
)

func init() {
	app.RegisterPlugin(app.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	configs := make([]xxxljob.Config, 0)
	if err := config.UnmarshalKey("xxljob", configs); err != nil {
		return err
	}
	if len(configs) == 0 {
		configs = append(configs, *xxxljob.DefaultConfig())
	}
	for _, config := range configs {
		cli := xxxljob.New(&config)
		cli.Init()
		set(config.Name, cli)
	}
	return nil
}
