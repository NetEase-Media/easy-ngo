package pluginxkafka

import (
	"sync"

	"github.com/NetEase-Media/easy-ngo/clients/xkafka"
)

var (
	mu          sync.RWMutex
	consumerMap = make(map[string]*xkafka.Consumer)
	producerMap = make(map[string]*xkafka.Producer)
)

func GetConsumerByKey(name string) *xkafka.Consumer {
	mu.RLock()
	defer mu.RUnlock()
	return consumerMap[name]
}

func GetProducerByKey(name string) *xkafka.Producer {
	mu.RLock()
	defer mu.RUnlock()
	return producerMap[name]
}

func GetConsumer() *xkafka.Consumer {
	return GetConsumerByKey("default")
}

func GetProducer(name string) *xkafka.Producer {
	return GetProducerByKey("default")
}

func setConsumer(name string, consumer *xkafka.Consumer) {
	mu.Lock()
	defer mu.Unlock()
	consumerMap[name] = consumer
}

func setProducer(name string, producer *xkafka.Producer) {
	mu.Lock()
	defer mu.Unlock()
	producerMap[name] = producer
}
