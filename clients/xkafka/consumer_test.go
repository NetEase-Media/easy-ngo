// Copyright 2022 NetEase Media Technology（Beijing）Co., Ltd.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xkafka

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/IBM/sarama"

	"github.com/stretchr/testify/assert"
)

const (
	// KAFKAVERSION = defaultVersion
	// KAFKAADDR    = "kafka:9092"
	KAFKAADDR        = "localhost:9092"
	KAFKAVERSION     = "0.11.0.0"
	KAFKATOPICNORMAL = "test"
	KAFKATOPICPREX   = "ngo-kafka-test-"
)

// 2个消费者消费，数据不会重复
func TestConsumer_TwoConsumer(t *testing.T) {
	topic := KAFKATOPICPREX + strconv.Itoa(int(time.Now().UnixNano()/1e6))
	opts := DefaultConfig()
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	opts.Consumer.Group = "ngo"
	c, err := NewConsumer(opts)
	assert.NoError(t, err)
	c2, err := NewConsumer(opts)
	var c1Count int64 = 0
	var c2Count int64 = 0
	c.AddListener(topic, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		atomic.AddInt64(&c1Count, 1)
		//log.Info("c1:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
		fmt.Print("c1:" + message.Value + ":" + strconv.Itoa(int(message.Partition)) + ",offset:" + strconv.Itoa(int(message.Offset)))
	}})
	c2.AddListener(topic, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		atomic.AddInt64(&c2Count, 1)
		//log.Info("c2:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
		fmt.Print("c2:" + message.Value + ":" + strconv.Itoa(int(message.Partition)) + ",offset:" + strconv.Itoa(int(message.Offset)))
	}})
	go c.Start()
	go c2.Start()
	time.Sleep(5 * time.Second)
	p, err := NewProducer(opts)
	assert.NoError(t, err)
	defer p.Close()
	i := 0
	for i < 100 {
		p.Send(topic, strconv.Itoa(i), nil)
		i++
	}
	time.Sleep(20 * time.Second)
	assert.Equal(t, int64(100), atomic.LoadInt64(&c1Count)+atomic.LoadInt64(&c2Count))
	c.Stop()
	c2.Stop()
}

// 1个消费者，不自动提交commit，手动提交
func TestConsumerEnabledAutoCommit(t *testing.T) {
	topic := KAFKATOPICPREX + "autocommit-" + strconv.Itoa(int(time.Now().UnixNano()/1e6))
	t.Logf("topic name: %s", topic)
	opts := DefaultConfig()
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	opts.Consumer.Group = "ngo"
	//手动提交commit
	opts.Consumer.EnableAutoCommit = false
	c, err := NewConsumer(opts)
	assert.NoError(t, err)
	var c1Count int64 = 0
	c.AddListener(topic, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		atomic.AddInt64(&c1Count, 1)
		//log.Info("c1:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
		ack.Acknowledge()
	}})
	c.Start()
	//生产消息
	p, err := NewProducer(opts)
	assert.NoError(t, err)
	defer p.Close()
	i := 0
	for i < 10 {
		p.Send(topic, strconv.Itoa(i), nil)
		i++
	}
	time.Sleep(5 * time.Second)
	assert.Equal(t, int64(10), atomic.LoadInt64(&c1Count))
	c.Stop()
	//再次进行消费
	c2, err := NewConsumer(opts)
	assert.NoError(t, err)
	var c2Count int64 = 0
	c2.AddListener(topic, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		atomic.AddInt64(&c2Count, 1)
		//log.Info("c2:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
		ack.Acknowledge()
	}})
	c2.Start()
	time.Sleep(5 * time.Second)
	c2.Stop()
	assert.Equal(t, int64(0), atomic.LoadInt64(&c2Count))
}

// 1个消费者，不自动提交commit且不手动提交
func TestConsumerNotEnabledAutoCommit(t *testing.T) {
	topic := KAFKATOPICPREX + "autocommit-" + strconv.Itoa(int(time.Now().UnixNano()/1e6))
	t.Logf("topic name: %s", topic)
	opts := DefaultConfig()
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	//消费后不手动提交
	//log.Info("opts.Consumer.EnableAutoCommit = false，模拟消费不提交")
	opts.Consumer.InitialOffset = sarama.OffsetOldest
	opts.Consumer.EnableAutoCommit = false
	opts.Consumer.Group = "ngo"
	c1, err := NewConsumer(opts)
	assert.NoError(t, err)
	c1.AddListener(topic, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		//log.Info("c1:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
	}})
	c1.Start()
	//生产5条消息
	p, err := NewProducer(opts)
	assert.NoError(t, err)
	defer p.Close()
	total := 5
	i := 0
	for i < total {
		p.Send(topic, strconv.Itoa(i), nil)
		i++
	}
	time.Sleep(5 * time.Second)
	c1.Stop()
	//消费后手动提交
	//log.Info("opts.Consumer.EnableAutoCommit = false，模拟手动提交")
	c2, err := NewConsumer(opts)
	assert.NoError(t, err)
	var c2Count int64 = 0
	c2.AddListener(topic, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		atomic.AddInt64(&c2Count, 1)
		//log.Info("c2:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
		ack.Acknowledge()
	}})
	c2.Start()
	time.Sleep(5 * time.Second)
	assert.Equal(t, int64(total), atomic.LoadInt64(&c2Count))
	c2.Stop()
	//再次进行消费
	//log.Info("验证队列中无可消费的消息")
	c3, err := NewConsumer(opts)
	assert.NoError(t, err)
	var c3Count int64 = 0
	c3.AddListener(topic, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		atomic.AddInt64(&c3Count, 1)
		//log.Info("c3:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
		ack.Acknowledge()
	}})
	c3.Start()
	time.Sleep(5 * time.Second)
	c3.Stop()
	assert.Equal(t, int64(0), atomic.LoadInt64(&c3Count))
}

// 1个消费者，消费不同group下的消息
func TestConsumerTwoGroup(t *testing.T) {
	topic := KAFKATOPICPREX + strconv.Itoa(int(time.Now().UnixNano()/1e6))
	t.Logf("topic name: %s", topic)
	opts := DefaultConfig()
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	opts.Consumer.InitialOffset = sarama.OffsetOldest
	opts.Consumer.Group = "test1"
	c, err := NewConsumer(opts)
	assert.NoError(t, err)
	var c1Count int64 = 0
	c.AddListener(topic, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		atomic.AddInt64(&c1Count, 1)
		//log.Info("c1:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
	}})
	c.Start()
	p, err := NewProducer(opts)
	assert.NoError(t, err)
	defer p.Close()
	i := 0
	for i < 5 {
		p.Send(topic, strconv.Itoa(i), nil)
		i++
	}
	time.Sleep(3 * time.Second)
	assert.Equal(t, int64(5), atomic.LoadInt64(&c1Count))
	c.Stop()
	timeStr := "test2"
	opts.Consumer.Group = timeStr
	c1, err := NewConsumer(opts)
	var c2Count int64 = 0
	c1.AddListener(topic, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		atomic.AddInt64(&c2Count, 1)
		//log.Info("c2:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
	}})
	c1.Start()
	time.Sleep(3 * time.Second)
	assert.Equal(t, int64(5), atomic.LoadInt64(&c2Count))
	c1.Stop()
	//同组再次消费
	opts.Consumer.Group = timeStr
	c2, err := NewConsumer(opts)
	var c3Count int64 = 0
	c2.AddListener(topic, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		atomic.AddInt64(&c3Count, 1)
		//log.Info("c3:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
	}})
	c2.Start()
	time.Sleep(3 * time.Second)
	assert.Equal(t, int64(0), atomic.LoadInt64(&c3Count))
	c2.Stop()
}

// 添加监听，topic为空
func TestConsumerAddListenerTopicNil(t *testing.T) {
	opts := DefaultConfig()
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	opts.Consumer.Group = "ngo"
	c0, err := NewConsumer(opts)
	assert.NoError(t, err)
	assert.Panics(t, func() {
		c0.AddListener("", &listener{func(message ConsumerMessage, ack *Acknowledgment) {
			//log.Info("c0:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
		}})
	})
}

// AddListener listener=nil
func TestConsumerAddListenerNil(t *testing.T) {
	opts := DefaultConfig()
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	opts.Consumer.Group = "ngo"
	c0, err := NewConsumer(opts)
	assert.NoError(t, err)
	assert.Panics(t, func() {
		c0.AddListener(KAFKATOPICNORMAL, nil)
	}, "topic must not be empty")
}

// 启动panic验证
func TestConsumerStart(t *testing.T) {
	opts := DefaultConfig()
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	opts.Consumer.Group = "ngo"
	c0, err := NewConsumer(opts)
	assert.NoError(t, err)
	assert.Panics(t, func() {
		c0.AddListener(KAFKATOPICNORMAL, nil)
	}, "listener must not be nil")
	assert.Panics(t, func() { c0.Start() }, "empty topic listener")
}

func TestConsumerStartDup(t *testing.T) {
	opts := DefaultConfig()
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	opts.Consumer.Group = "ngo"
	c0, err := NewConsumer(opts)
	assert.NoError(t, err)
	c0.AddListener(KAFKATOPICNORMAL, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		//log.Info("c0:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
	}})
	c0.Start()
	assert.Panics(t, func() {
		c0.Start()
	}, "duplicated start")
	c0.Stop()
}

func TestConsumerStopException(t *testing.T) {
	opts := DefaultConfig()
	opts.Addr = []string{KAFKAADDR}
	opts.Version = "0.2..3"
	opts.Consumer.Group = "ngo"
	//初始化失败
	c0, _ := NewConsumer(opts)
	assert.Panics(t, func() {
		c0.Stop()
	}, "runtime error: invalid memory address or nil pointer dereference")
}

func TestConsumerStopCancelNotNil(t *testing.T) {
	opts := DefaultConfig()
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	opts.Consumer.Group = "ngo"
	c0, err := NewConsumer(opts)
	assert.NoError(t, err)
	c0.AddListener(KAFKATOPICNORMAL, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		//log.Info("c0:" + message.Value + "：" + strconv.Itoa(int(message.Partition)) + "，offset:" + strconv.Itoa(int(message.Offset)))
	}})
	c0.Start()
	c0.Stop()
}

func TestConsumerListenerErr(t *testing.T) {
	topic := KAFKATOPICPREX + strconv.Itoa(int(time.Now().UnixNano()/1e6))
	println(topic)
	opts := DefaultConfig()
	opts.Addr = []string{KAFKAADDR}
	opts.Version = KAFKAVERSION
	opts.Consumer.Group = "ngo"
	c, err := NewConsumer(opts)
	assert.NoError(t, err)
	var i int32 = 0
	c.AddListener(topic, &listener{func(message ConsumerMessage, ack *Acknowledgment) {
		atomic.AddInt32(&i, 1)
		println(message.Partition, message.Offset, message.Value)
		if atomic.LoadInt32(&i) > 40 {
			panic(fmt.Sprintf("panic: %d %d %s", message.Partition, message.Offset, message.Value))
		}
	}})
	go c.Start()
	time.Sleep(5 * time.Second)
	p, err := NewProducer(opts)
	assert.NoError(t, err)
	defer p.Close()
	for j := 0; j < 50; j++ {
		p.Send(topic, strconv.Itoa(j), nil)
	}
	time.Sleep(10 * time.Second)
	c.Stop()
}

type listener struct {
	listenFn func(message ConsumerMessage, ack *Acknowledgment)
}

func (l *listener) Listen(message ConsumerMessage, ack *Acknowledgment) {
	l.listenFn(message, ack)
}
