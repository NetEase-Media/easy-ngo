package pluginxmemcache

import (
	"context"

	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/NetEase-Media/easy-ngo/clients/xmemcache"
	"github.com/NetEase-Media/easy-ngo/config"
)

func init() {
	app.RegisterPlugin(app.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	configs := make([]xmemcache.Config, 0)
	if err := config.UnmarshalKey("memcache", configs); err != nil {
		return err
	}
	if len(configs) == 0 {
		configs = append(configs, *xmemcache.DefaultConfig())
	}
	for _, config := range configs {
		cli, err := xmemcache.New(&config)
		if err != nil {
			return err
		}
		set(config.Name, cli)
	}
	return nil
}
