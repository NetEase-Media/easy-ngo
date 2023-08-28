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
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/NetEase-Media/easy-ngo/xlog"
)

type ConsumerMessage struct {
	Topic     string
	Key       string
	Value     string
	Partition int32
	Offset    int64
}

type Listener interface {
	Listen(ConsumerMessage, *Acknowledgment)
}

type BatchListener interface {
	Listen([]ConsumerMessage, *Acknowledgment)
	BatchCount() int
}

// Consumer 是一个group的消费者
type Consumer struct {
	client         sarama.ConsumerGroup
	config         Config
	ctx            context.Context
	cancel         func()
	runChan        chan struct{}
	listeners      map[string]Listener
	batchListeners map[string]BatchListener
}

func (c *Consumer) AddListener(topic string, listener Listener) {
	if len(topic) == 0 {
		panic("topic must not be empty")
	}
	if listener == nil {
		panic("listener must not be nil")
	}
	c.listeners[topic] = listener
}

func (c *Consumer) AddBatchListener(topic string, listener BatchListener) {
	if len(topic) == 0 {
		panic("topic must not be empty")
	}
	if listener == nil {
		panic("listener must not be nil")
	}
	c.batchListeners[topic] = listener
}

// Start 启动后台消费任务
func (c *Consumer) Start() {
	if len(c.listeners) == 0 && len(c.batchListeners) == 0 {
		panic("empty topic listener")
	}

	// 当前不允许多个后台消费任务
	if c.ctx != nil {
		panic("duplicated start")
	}

	h := &consumerHandler{
		consumer: c,
		ready:    make(chan struct{}),
		config:   &c.config,
	}

	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.runChan = make(chan struct{})
	tmap := make(map[string]struct{}, len(c.listeners)+len(c.batchListeners))
	for k := range c.listeners {
		tmap[k] = struct{}{}
	}
	for k := range c.batchListeners {
		tmap[k] = struct{}{}
	}

	topics := make([]string, 0, len(tmap))
	for k := range tmap {
		topics = append(topics, k)
	}

	go func() {
		defer close(c.runChan)
		for {
			// 当服务的rebalance后会返回
			if err := c.client.Consume(c.ctx, topics, h); err != nil {
				xlog.Errorf("kafka consume failed: %s", err.Error())
				time.Sleep(time.Millisecond * 200) // 睡眠防止异常之后死循环占满CPU
			}

			if c.ctx.Err() != nil {
				return
			}

			select {
			case <-h.ready:
				h.ready = make(chan struct{})
			default:
			}
		}
	}()
	<-h.ready
	xlog.Infof("consumer up and running")
}

// Stop 停止后台消费任务
func (c *Consumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
		<-c.runChan
	}
	return c.client.Close()
}

func NewConsumer(conf *Config) (*Consumer, error) {
	config, err := newConsumerConfig(conf)
	if err != nil {
		return nil, err
	}
	c, err := sarama.NewConsumerGroup(conf.Addr, conf.Consumer.Group, config)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		client:         c,
		config:         *conf,
		listeners:      make(map[string]Listener, 8),
		batchListeners: make(map[string]BatchListener, 8),
	}, nil
}

func newConsumerConfig(opt *Config) (*sarama.Config, error) {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(opt.Version)
	if err != nil {
		return nil, err
	}
	config.Version = version

	config.Metadata.RefreshFrequency = time.Second * 10

	config.Net.MaxOpenRequests = opt.MaxOpenRequests
	config.Net.DialTimeout = opt.DialTimeout
	config.Net.ReadTimeout = opt.ReadTimeout
	config.Net.WriteTimeout = opt.WriteTimeout

	config.Net.SASL.Enable = opt.SASL.Enable
	config.Net.SASL.Mechanism = opt.SASL.Mechanism
	config.Net.SASL.User = opt.SASL.User
	config.Net.SASL.Password = opt.SASL.Password
	config.Net.SASL.Handshake = opt.SASL.Handshake

	config.Metadata.Retry.Max = opt.Metadata.Retries
	config.Metadata.Timeout = opt.Metadata.Timeout

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = opt.Consumer.EnableAutoCommit
	config.Consumer.Offsets.AutoCommit.Interval = opt.Consumer.AutoCommitInterval
	config.Consumer.Offsets.Initial = opt.Consumer.InitialOffset
	config.Consumer.Offsets.Retry.Max = opt.Consumer.Retries
	config.Consumer.Group.Session.Timeout = opt.Consumer.SessionTimeout
	config.Consumer.Fetch.Min = opt.Consumer.MinFetchBytes
	config.Consumer.Fetch.Default = opt.Consumer.DefaultFetchBytes
	config.Consumer.Fetch.Max = opt.Consumer.MaxFetchBytes
	config.Consumer.MaxWaitTime = opt.Consumer.MaxFetchWait
	return config, nil
}

// consumerHandler 用来运行消费者后台任务
type consumerHandler struct {
	consumer     *Consumer
	ready        chan struct{}
	config       *Config
	enable       bool
	enableTracer bool
}

// Setup 在启动前执行
func (ch *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	close(ch.ready)
	return nil
}

// Cleanup 在结束后执行
func (ch *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 在循环中消费message
func (ch *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	topic := claim.Topic()
	if batchListener := ch.consumer.batchListeners[topic]; batchListener != nil {
		count := batchListener.BatchCount()
		if count < 1 {
			count = 1
		}
		msgArr := make([]*sarama.ConsumerMessage, 0, count)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case message := <-claim.Messages():
				if message == nil {
					xlog.Infof("channel is closed")
					return nil
				}
				msgArr = append(msgArr, message)
				if len(msgArr) >= count {
					ch.batchListen(batchListener, session, msgArr)
					msgArr = make([]*sarama.ConsumerMessage, 0, count)
				}
			case <-ticker.C:
				if len(msgArr) > 0 {
					ch.batchListen(batchListener, session, msgArr)
					msgArr = make([]*sarama.ConsumerMessage, 0, count)
				}
			}
		}
	}
	return nil
}

func (ch *consumerHandler) listen(listener Listener, session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) {
	msg := ConsumerMessage{
		Topic:     message.Topic,
		Key:       string(message.Key),
		Value:     string(message.Value),
		Partition: message.Partition,
		Offset:    message.Offset,
	}
	ack := &Acknowledgment{
		ch:      ch,
		session: session,
		message: message,
	}
	begin := time.Now()
	defer func() {
		var err error
		switch r := recover().(type) {
		case nil:
		case error:
			err = r
		default:
			err = fmt.Errorf("unexpected panic value: %#v", r)
		}
		if err != nil {
			json, _ := json.Marshal(&msg)
			xlog.Errorf("consumer handle error: %v, message: %s", err, json)
		}
		ch.collect(message.Topic, message.Partition, len(message.Value), time.Since(begin), err)
	}()

	listener.Listen(msg, ack)
	// if auto commit, mark message
	if ch.consumer.config.Consumer.EnableAutoCommit {
		session.MarkMessage(message, "")
	}
}

func (ch *consumerHandler) batchListen(listener BatchListener, session sarama.ConsumerGroupSession, msgArr []*sarama.ConsumerMessage) {
	begin := time.Now()
	topic := msgArr[0].Topic
	partition := msgArr[0].Partition
	var msgBytes int
	defer func() {
		var err error
		switch r := recover().(type) {
		case nil:
		case error:
			err = r
		default:
			err = fmt.Errorf("unexpected panic value: %#v", r)
		}
		if err != nil {
			xlog.Errorf("batch consumer handle error: %v", err)
		}
		ch.collect(topic, partition, msgBytes, time.Since(begin), err)
	}()
	msgs := make([]ConsumerMessage, 0, len(msgArr))
	for _, message := range msgArr {
		msgs = append(msgs, ConsumerMessage{
			Topic:     message.Topic,
			Key:       string(message.Key),
			Value:     string(message.Value),
			Partition: message.Partition,
			Offset:    message.Offset,
		})
		msgBytes += len(message.Value)
	}
	ack := &Acknowledgment{
		ch:      ch,
		session: session,
		message: msgArr[len(msgArr)-1],
	}
	listener.Listen(msgs, ack)

	// if auto commit, mark message
	if ch.consumer.config.Consumer.EnableAutoCommit {
		session.MarkMessage(msgArr[len(msgArr)-1], "")
	}
}

// collect 生成监控数据发送到收集器
func (ch *consumerHandler) collect(topic string, partition int32, msgBytes int, cost time.Duration, err error) {
	// if !metrics.IsMetricsEnabled() {
	// 	return
	// }
	// r := &kafka.StatsRecord{
	// 	Broker:       ch.opt.Addr[0],
	// 	Topic:        topic,
	// 	Partition:    partition,
	// 	Cost:         cost,
	// 	MessageBytes: int64(msgBytes),
	// 	Err:          err,
	// }
	// collectors.KafkaCollector().OnReceived(r)
	if !ch.enable {
		return
	}
}

type Acknowledgment struct {
	ch      *consumerHandler
	session sarama.ConsumerGroupSession
	message *sarama.ConsumerMessage
}

func (a *Acknowledgment) Acknowledge() {
	if !a.ch.consumer.config.Consumer.EnableAutoCommit {
		a.session.MarkMessage(a.message, "")
		a.session.Commit()
	}
}
