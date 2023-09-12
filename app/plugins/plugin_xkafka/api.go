package pluginxkafka

import (
	"sync"

	"github.com/NetEase-Media/easy-ngo/clients/xkafka"
)

var (
	mu sync.RWMutex
)

func GetConsumerByKey(name string) *xkafka.Consumer {
	mu.RLock()
	defer mu.RUnlock()
	return consumers[name]
}

func GetProducerByKey(name string) *xkafka.Producer {
	mu.RLock()
	defer mu.RUnlock()
	return producers[name]
}

func GetConsumer() *xkafka.Consumer {
	return GetConsumerByKey("default")
}

func GetProducer() *xkafka.Producer {
	return GetProducerByKey("default")
}
