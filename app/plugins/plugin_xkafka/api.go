// Copyright 2022 NetEase Media Technology（Beijing）Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
