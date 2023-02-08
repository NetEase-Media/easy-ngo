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

package xotel

import (
	"github.com/NetEase-Media/easy-ngo/observability/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	resource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type Option struct {
	// 服务名
	ServiceName string
	// collector地址
	Endpoint string
	// 支持jaeger,otel,stdout
	// 默认 stdout
	ExporterName string
	// sample 采样率 0 < x<=1
	SamplerRatio float64
}

func DefaultOption() *Option {
	return &Option{
		ServiceName:  "localhost",
		Endpoint:     "",
		ExporterName: "stdout",
		SamplerRatio: 0.1,
	}
}

func New(option *Option) tracing.Provider {
	var exp sdktrace.SpanExporter
	switch option.ExporterName {
	case "stdout":
		exp = newStdoutExporter()
	case "jaeger":
		exp = newJaegerExporter(option.Endpoint)
	default:
		exp = newStdoutExporter()
	}
	return NewProvider(option, exp)
}

// TODO: 后续提供 resource sampler idGenerator spanlimits的可定制化
// SpanExporter 将otel生成的span数据做转化并上报到相应的collector，比如zipkin jaeger
// Sampler 采样器
// IDGenerator spanid traceid生成器
// SpanLimits 限制event attribute的数量
func NewProvider(option *Option, exp sdktrace.SpanExporter) tracing.Provider {
	res := resource.NewSchemaless(
		semconv.TelemetrySDKLanguageGo,
		semconv.ServiceNameKey.String(option.ServiceName),
	)
	provider := sdktrace.NewTracerProvider(
		// 设置导出exporter
		sdktrace.WithBatcher(exp),
		// 设置采样器
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(option.SamplerRatio))),
		sdktrace.WithResource(res),
	)
	setGlobalProvider(provider)
	return provider
}

func setGlobalProvider(p tracing.Provider) {
	// 设置全局的provider
	otel.SetTracerProvider(p)
	// 设置全局的propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

// 默认的otel provider实现
func DefaultProvider() tracing.Provider {
	p := New(DefaultOption())
	// 设置全局的provider
	otel.SetTracerProvider(p)
	// 设置全局的propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return p
}

// exporters
// jaeger exporter
func newJaegerExporter(url string) sdktrace.SpanExporter {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic("jaeger url error")
	}
	return exp
}

// stdout exporter
func newStdoutExporter() sdktrace.SpanExporter {
	exp, err := stdouttrace.New()
	if err != nil {
		panic(err)
	}
	return exp
}
