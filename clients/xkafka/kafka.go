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

const (
	defaultVersion = "2.1.0"
)

func New(opt *Config) (*Kafka, error) {

	k := &Kafka{
		Opt: opt,
	}
	hasProducer := true // TODO default must have a producer
	hasConsumer := opt.Consumer.Group != ""
	if hasProducer {
		p, err := NewProducer(opt)
		if err != nil {
			return nil, err
		}
		k.Producer = p
	}
	if hasConsumer {
		co, err := NewConsumer(opt)
		if err != nil {
			return nil, err
		}
		k.Consumer = co
	}
	return k, nil
}

type Kafka struct {
	Opt      *Config
	Consumer *Consumer
	Producer *Producer
}
