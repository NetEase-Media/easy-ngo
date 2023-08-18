// // Copyright 2022 NetEase Media Technology（Beijing）Co., Ltd.
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// // 	http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.

package xkafka

// import (
// 	"context"
// 	"strconv"
// 	"sync"
// 	"testing"

// 	"github.com/NetEase-Media/easy-ngo/xlog/xfmt"
// 	"github.com/stretchr/testify/assert"
// )

// func TestProducerSend(t *testing.T) {
// 	opts := NewDefaultOptions()
// 	opts.Addr = []string{KAFKAADDR}
// 	opts.Version = KAFKAVERSION
// 	p, err := NewProducer(opts, &xfmt.XFmt{}, nil, nil)
// 	assert.NoError(t, err)
// 	defer p.Close()

// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	//异步发送消息
// 	p.Send("ngo-kafka-test-produce", "test-ahah-Send", func(meta *RecordMetadata, err error) {
// 		assert.NoError(t, err)
// 		wg.Done()
// 	})
// 	//同步发送消息
// 	err = p.SyncSend(context.Background(), "ngo-kafka-test-produce", "test-ahah-SyncSend")
// 	assert.NoError(t, err)
// 	wg.Wait()
// }

// // 版本数据异常
// func TestProducerSend_KafkaVersion_Exception(t *testing.T) {
// 	opts := NewDefaultOptions()
// 	opts.Addr = []string{KAFKAADDR}
// 	opts.Version = "1.0"
// 	_, err := NewProducer(opts, &xfmt.XFmt{}, nil, nil)
// 	assert.Error(t, err)
// 	opts.Version = "0.0.a.0"
// 	_, err1 := NewProducer(opts, &xfmt.XFmt{}, nil, nil)
// 	assert.Error(t, err1)
// 	opts.Version = "1.0.a.0"
// 	_, err2 := NewProducer(opts, &xfmt.XFmt{}, nil, nil)
// 	assert.Error(t, err2)
// }

// func TestCloseProducerSend(t *testing.T) {
// 	opts := NewDefaultOptions()
// 	opts.Addr = []string{KAFKAADDR}
// 	opts.Version = KAFKAVERSION
// 	p, _ := NewProducer(opts, &xfmt.XFmt{}, nil, nil)
// 	p.SyncSend(context.Background(), "ngo-kafka-test-produce", "test-ahah-SyncSend")
// 	p.Close()
// 	//同步发送消息
// 	assert.Panics(t, func() {
// 		p.SyncSend(context.Background(), "ngo-kafka-test-produce", "test-ahah-SyncSend")
// 	})
// }

// //单测无法模拟，异常case条件允许走功能测试，发消息过程中关闭kafka
// /*func TestProducerSend_SendTimeout(t *testing.T) {
// 	opts := NewDefaultOptions()
// 	opts.Addr = []string{"127.0.0.1:9092"}
// 	opts.Version = "0.11.0.0"
// 	p, err := NewProducer(opts)
// 	p.Close()
// 	//同步发送消息
// 	err = p.SyncSend("ngo-test", "test-ahah-SyncSend")
// 	assert.NoError(t, err)
// }*/

// func BenchmarkProducerSend(b *testing.B) {
// 	opts := NewDefaultOptions()
// 	opts.Addr = []string{"kafka:9092"}
// 	p, err := NewProducer(opts, &xfmt.XFmt{}, nil, nil)
// 	assert.NoError(b, err)
// 	p.run()
// 	defer p.Close()

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		p.Send("test", strconv.Itoa(i), func(meta *RecordMetadata, err error) {
// 			assert.NoError(b, err)
// 		})
// 	}
// 	b.StopTimer()
// }
