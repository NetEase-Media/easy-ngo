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

package xkafka

import (
	"sync"
)

var (
	mu          sync.RWMutex
	consumerMap = make(map[string]*Consumer)
	producerMap = make(map[string]*Producer)
)

func GetConsumer(name string) *Consumer {
	mu.RLock()
	defer mu.RUnlock()
	return consumerMap[name]
}

func GetProducer(name string) *Producer {
	mu.RLock()
	defer mu.RUnlock()
	return producerMap[name]
}

func SetConsumer(name string, consumer *Consumer) {
	mu.Lock()
	defer mu.Unlock()
	consumerMap[name] = consumer
}

func SetProducer(name string, producer *Producer) {
	mu.Lock()
	defer mu.Unlock()
	producerMap[name] = producer
}
