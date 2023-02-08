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

	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
	"github.com/Shopify/sarama"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// producer trace before
func traceSendBefore(ctx context.Context, m *sarama.ProducerMessage) tracer.Span {
	// start span
	tr := tracer.GetTracer("samara")
	propagator := tracer.GetTextMapPropagator()
	newCtx, span := tr.Start(
		ctx, "kafka-producer",
		tracer.WithAttributes(
			semconv.MessagingSystemKey.String("kafka"),
			semconv.MessagingOperationKey.String("SyncSend"),
			semconv.MessagingDestinationKey.String(m.Topic),
		),
		tracer.WithSpanKind(tracer.SpanKindProducer),
	)
	// 注入spanid信息
	propagator.Inject(newCtx, NewProducerMessageCarrier(m))
	return span
}

// producer trace after
func traceSendAfter(span tracer.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
	return
}

// producer consume before
func traceConsumeBefore(m *sarama.ConsumerMessage) tracer.Span {
	// start span
	tr := tracer.GetTracer("samara")
	propagator := tracer.GetTextMapPropagator()
	parentContext := propagator.Extract(context.Background(), NewConsumerMessageCarrier(m))
	_, span := tr.Start(
		parentContext, "kafka-consumer",
		tracer.WithAttributes(
			semconv.MessagingSystemKey.String("kafka"),
			semconv.MessagingOperationKey.String("Consumer"),
			semconv.MessagingDestinationKey.String(m.Topic),
			semconv.MessagingKafkaPartitionKey.Int64(int64(m.Partition)),
			semconv.MessagingMessageIDKey.Int64(m.Offset),
		),
		tracer.WithSpanKind(tracer.SpanKindConsumer),
	)
	return span
}

// producer trace after
func traceConsumeAfter(span tracer.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
	return
}

// producer message carrier
type ProducerMessageCarrier struct {
	*sarama.ProducerMessage
	kvs  map[string]string
	keys []string
}

func NewProducerMessageCarrier(m *sarama.ProducerMessage) *ProducerMessageCarrier {
	carrier := &ProducerMessageCarrier{
		ProducerMessage: m,
		kvs:             make(map[string]string, len(m.Headers)),
		keys:            make([]string, len(m.Headers)),
	}
	for _, header := range m.Headers {
		carrier.kvs[string(header.Key)] = string(header.Value)
		carrier.keys = append(carrier.keys, string(header.Key))
	}
	return carrier
}

func (carrier ProducerMessageCarrier) Get(key string) string {
	return carrier.kvs[key]
}

func (carrier ProducerMessageCarrier) Set(key string, value string) {
	// 存在替换，不存在添加
	if _, ok := carrier.kvs[key]; ok {
		for i := 0; i < len(carrier.keys); i++ {
			if string(carrier.Headers[i].Key) == key {
				carrier.Headers[i].Value = []byte(value)
			}
		}
	} else {
		carrier.Headers = append(carrier.Headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(value),
		})
	}
}

func (carrier ProducerMessageCarrier) Keys() []string {
	return carrier.keys
}

// consumer message carrier
type ConsumerMessageCarrier struct {
	*sarama.ConsumerMessage
	kvs  map[string]string
	keys []string
}

func NewConsumerMessageCarrier(m *sarama.ConsumerMessage) *ConsumerMessageCarrier {
	carrier := &ConsumerMessageCarrier{
		ConsumerMessage: m,
		kvs:             make(map[string]string, len(m.Headers)),
		keys:            make([]string, len(m.Headers)),
	}
	for _, header := range m.Headers {
		carrier.kvs[string(header.Key)] = string(header.Value)
		carrier.keys = append(carrier.keys, string(header.Key))
	}
	return carrier
}

func (carrier ConsumerMessageCarrier) Get(key string) string {
	return carrier.kvs[key]
}

func (carrier ConsumerMessageCarrier) Set(key string, value string) {
	// 存在替换，不存在添加
	if _, ok := carrier.kvs[key]; ok {
		for i := 0; i < len(carrier.keys); i++ {
			if string(carrier.Headers[i].Key) == key {
				carrier.Headers[i].Value = []byte(value)
			}
		}
	} else {
		carrier.Headers = append(carrier.Headers, &sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(value),
		})
	}
}

func (carrier ConsumerMessageCarrier) Keys() []string {
	return carrier.keys
}
