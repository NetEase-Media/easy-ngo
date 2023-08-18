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
// 	"errors"
// 	"sync"
// 	"time"

// 	"github.com/NetEase-Media/easy-ngo/observability/metrics"
// 	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
// 	"github.com/NetEase-Media/easy-ngo/xlog"
// 	"github.com/Shopify/sarama"
// 	"github.com/prometheus/client_golang/prometheus"
// )

// const (
// 	metricProducerReadTotalName         = "kafka_producer_read_total"
// 	metricProducerReadDurationName      = "kafka_producer_read_duration"
// 	metricProducerReadDurationRangeName = "kafka_producer_read_duration_range"
// 	metricProducerErrorName             = "kafka_producer_error"
// )

// var (
// 	metricProducerReadTotal         metrics.Counter
// 	metricProducerReadDuration      metrics.Gauge
// 	metricProducerReadDurationRange metrics.Histogram
// 	metricProducerError             metrics.Counter
// )

// type ProducerMessage struct {
// 	Topic string
// 	Key   string
// 	Value string
// }

// type RecordMetadata struct {
// 	Topic     string
// 	KeySize   int
// 	ValueSize int
// 	Offset    int64
// 	Partition int32
// }

// type Producer struct {
// 	client  sarama.AsyncProducer
// 	opt     Option
// 	logger  xlog.Logger
// 	runChan chan struct{}
// 	wg      sync.WaitGroup
// 	metrics metrics.Provider
// 	tracer  tracer.Provider
// }

// type Callback func(*RecordMetadata, error)

// func (p *Producer) Options() Option {
// 	return p.opt
// }

// // Send 是异步发送接口
// func (p *Producer) Send(topic, value string, cb Callback) {
// 	p.SendMessage(ProducerMessage{Topic: topic, Value: value}, cb)
// }

// // SendMessage 是异步发送接口
// func (p *Producer) SendMessage(message ProducerMessage, cb Callback) {

// 	var span tracer.Span

// 	meta := newMetaData()
// 	meta.cb = cb
// 	if p.opt.EnableTracer {
// 		meta.cb = func(rm *RecordMetadata, e error) {
// 			// 结束记录trace
// 			traceSendAfter(span, e)
// 			cb(rm, e)
// 		}
// 	}

// 	m := &sarama.ProducerMessage{
// 		Topic:    message.Topic,
// 		Key:      nil,
// 		Value:    sarama.StringEncoder(message.Value),
// 		Metadata: meta,
// 	}
// 	// 开始记录span
// 	if p.opt.EnableTracer {
// 		span = traceSendBefore(context.Background(), m)
// 	}

// 	if len(message.Key) != 0 {
// 		m.Key = sarama.StringEncoder(message.Key)
// 	}
// 	p.client.Input() <- m
// }

// // SyncSend 是同步发送接口。
// func (p *Producer) SyncSend(ctx context.Context, topic, value string) error {
// 	return p.SyncSendMessage(ctx, ProducerMessage{Topic: topic, Value: value})
// }

// // SyncSendMessage 是同步发送接口。
// func (p *Producer) SyncSendMessage(ctx context.Context, message ProducerMessage) error {
// 	var span tracer.Span

// 	meta := newMetaData()
// 	meta.resChan = make(chan error)
// 	if p.opt.EnableTracer {
// 		meta.cb = func(rm *RecordMetadata, err error) {
// 			traceSendAfter(span, err)
// 		}
// 	}

// 	m := &sarama.ProducerMessage{
// 		Topic:    message.Topic,
// 		Key:      nil,
// 		Value:    sarama.StringEncoder(message.Value),
// 		Metadata: meta,
// 	}
// 	if len(message.Key) != 0 {
// 		m.Key = sarama.StringEncoder(message.Key)
// 	}

// 	// 开始记录trace
// 	if p.opt.EnableTracer {
// 		span = traceSendBefore(ctx, m)
// 	}

// 	p.client.Input() <- m

// 	timer := time.NewTimer(time.Second * 10)
// 	var rerr error
// 	defer timer.Stop()
// 	select {
// 	case err := <-meta.resChan:
// 		rerr = err
// 	case <-timer.C:
// 		// 防止异常卡死
// 		rerr = errors.New("send timeout")
// 	}

// 	return rerr
// }

// // run 启动后台任务，接收结果和错误
// func (p *Producer) run() {
// 	p.wg.Add(1)
// 	go p.receiveSuccess()

// 	p.wg.Add(1)
// 	go p.receiveError()
// }

// // receiveSuccess 接收成功回复，记录结果
// func (p *Producer) receiveSuccess() {
// 	defer p.wg.Done()
// 	for {
// 		select {
// 		case s, ok := <-p.client.Successes():
// 			if !ok {
// 				return
// 			}
// 			p.handle(s, nil)
// 		}
// 	}
// }

// // receiveError 接收错误回复
// func (p *Producer) receiveError() {
// 	defer p.wg.Done()
// 	for {
// 		select {
// 		case e, ok := <-p.client.Errors():
// 			if !ok {
// 				return
// 			}
// 			p.handle(e.Msg, e.Err)
// 		}
// 	}
// }

// // handle 处理异步消息的发送结果
// func (p *Producer) handle(msg *sarama.ProducerMessage, err error) {
// 	p.logger.Debugf("receive send response %+v, error %v", msg, err)
// 	meta := msg.Metadata.(*metaData)
// 	if meta.resChan != nil {
// 		meta.resChan <- err
// 		close(meta.resChan)
// 	}

// 	var ks, vs int
// 	if msg.Key != nil {
// 		ks = msg.Key.Length()
// 	}
// 	if msg.Value != nil {
// 		vs = msg.Value.Length()
// 	}

// 	if meta.cb != nil {
// 		meta.cb(&RecordMetadata{
// 			Topic:     msg.Topic,
// 			KeySize:   ks,
// 			ValueSize: vs,
// 			Offset:    msg.Offset,
// 			Partition: msg.Partition,
// 		}, err)
// 	}

// 	// if metrics.IsMetricsEnabled() {
// 	// 	r := &kafka.StatsRecord{
// 	// 		Broker:       p.opt.Addr[0],
// 	// 		Topic:        msg.Topic,
// 	// 		Partition:    msg.Partition,
// 	// 		Cost:         time.Since(meta.startTime),
// 	// 		MessageBytes: int64(vs),
// 	// 		Err:          err,
// 	// 	}
// 	// 	collectors.KafkaCollector().OnSent(r)
// 	// }
// 	p.collect(msg.Topic, msg.Partition, time.Since(meta.startTime), err)

// }

// func (p *Producer) collect(topic string, partition int32, cost time.Duration, err error) {
// 	// if metrics.IsMetricsEnabled() {
// 	// 	r := &kafka.StatsRecord{
// 	// 		Broker:       p.opt.Addr[0],
// 	// 		Topic:        msg.Topic,
// 	// 		Partition:    msg.Partition,
// 	// 		Cost:         time.Since(meta.startTime),
// 	// 		MessageBytes: int64(vs),
// 	// 		Err:          err,
// 	// 	}
// 	// 	collectors.KafkaCollector().OnSent(r)
// 	// }
// 	if p.metrics == nil {
// 		return
// 	}

// 	metricProducerReadTotal.With("host", p.opt.Addr[0], "topic", topic).Add(1)
// 	metricProducerReadDuration.With("host", p.opt.Addr[0], "topic", topic).Set(float64(cost) / 1e6)
// 	metricProducerReadDurationRange.With("host", p.opt.Addr[0], "topic", topic).Observe(float64(cost) / 1e6)
// 	if err != nil {
// 		metricProducerError.With("host", p.opt.Addr[0], "topic", topic).Add(1)
// 	}
// }

// // Close 关闭客户端，等待缓冲区完成读写再返回
// func (p *Producer) Close() {
// 	p.client.AsyncClose()
// 	p.wg.Wait()
// }

// // NewProducer 创建一个异步的生产者
// func NewProducer(opt *Option, logger xlog.Logger, metrics metrics.Provider, tracer tracer.Provider) (*Producer, error) {
// 	config, err := newProducerConfig(opt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	p, err := sarama.NewAsyncProducer(opt.Addr, config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	producer := &Producer{
// 		client:  p,
// 		opt:     *opt,
// 		runChan: make(chan struct{}),
// 		logger:  logger,
// 		metrics: metrics,
// 		tracer:  tracer,
// 	}

// 	if producer.metrics != nil {
// 		metricProducerReadTotal = producer.metrics.NewCounter(metricProducerReadTotalName, labelValues...)
// 		metricProducerReadDuration = producer.metrics.NewGauge(metricProducerReadDurationName, labelValues...)
// 		bukets := prometheus.ExponentialBuckets(.001, 10, 5)
// 		metricProducerReadDurationRange = producer.metrics.NewHistogram(metricProducerReadDurationRangeName, bukets, labelValues...)
// 		metricProducerError = producer.metrics.NewCounter(metricProducerErrorName, labelValues...)
// 	}

// 	producer.run() // 放到之后运行
// 	return producer, nil
// }

// func newProducerConfig(opt *Option) (*sarama.Config, error) {
// 	config := sarama.NewConfig()
// 	version, err := sarama.ParseKafkaVersion(opt.Version)
// 	if err != nil {
// 		return nil, err
// 	}
// 	config.Version = version
// 	config.ChannelBufferSize = 1024

// 	config.Net.MaxOpenRequests = opt.MaxOpenRequests
// 	config.Net.DialTimeout = opt.DialTimeout
// 	config.Net.ReadTimeout = opt.ReadTimeout
// 	config.Net.WriteTimeout = opt.WriteTimeout

// 	config.Net.SASL.Enable = opt.SASL.Enable
// 	config.Net.SASL.Mechanism = opt.SASL.Mechanism
// 	config.Net.SASL.User = opt.SASL.User
// 	config.Net.SASL.Password = opt.SASL.Password
// 	config.Net.SASL.Handshake = opt.SASL.Handshake

// 	config.Metadata.Retry.Max = opt.Metadata.Retries
// 	config.Metadata.Timeout = opt.Metadata.Timeout

// 	config.Producer.Return.Successes = true
// 	config.Producer.RequiredAcks = opt.Producer.Acks
// 	config.Producer.Timeout = opt.Producer.Timeout
// 	config.Producer.Retry.Max = opt.Producer.Retries
// 	config.Producer.Flush.Bytes = opt.Producer.MaxFlushBytes
// 	config.Producer.Flush.Messages = opt.Producer.MaxFlushMessages
// 	config.Producer.Flush.Frequency = opt.Producer.FlushFrequency
// 	config.Producer.Idempotent = opt.Producer.Idempotent

// 	return config, nil
// }

// func newMetaData() *metaData {
// 	return &metaData{
// 		startTime: time.Now(),
// 	}
// }

// // metaData 注册到message中，主要用来监控
// type metaData struct {
// 	startTime time.Time
// 	resChan   chan error // 回复channel，可以将异步调用变成同步
// 	cb        Callback
// }
