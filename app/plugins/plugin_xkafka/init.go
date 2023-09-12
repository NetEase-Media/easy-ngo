package pluginxkafka

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/NetEase-Media/easy-ngo/app"
	"github.com/NetEase-Media/easy-ngo/clients/xkafka"
	"github.com/NetEase-Media/easy-ngo/config"
)

var (
	producers = make(map[string]*xkafka.Producer, 1)
	consumers = make(map[string]*xkafka.Consumer, 1)
)

func init() {
	app.RegisterPlugin(app.Initialize, Initialize)
}

func Initialize(ctx context.Context) error {
	configs := make([]xkafka.Config, 0)
	if err := config.UnmarshalKey("kafka", &configs); err != nil {
		return err
	}
	if len(configs) == 0 {
		return errors.New("kafka config is empty")
	}
	for _, opt := range configs {
		opt.Consumer.InitialOffset = sarama.OffsetNewest
		cli, err := xkafka.New(&opt)
		if err != nil {
			panic("init kafka failed." + err.Error())
		}
		producers[opt.Name] = cli.Producer
		consumers[opt.Name] = cli.Consumer
	}
	return nil
}
