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

package xtracer

import (
	"context"
	"io"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.14.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 此处Tracer使用otel的API作为标准来实现
type Tracer = trace.Tracer

// 重要，此接口必须实现
// 依托实现的接口是TracerProvider
// 第三方实现的Tracer必须实现Provider接口
type Provider interface {
	trace.TracerProvider
	ForceFlush(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// otel/trace 包常量、类型、方法 快捷链接
// otel/trace 常量
const (
	FlagsSampled        = trace.FlagsSampled
	SpanKindUnspecified = trace.SpanKindUnspecified
	SpanKindInternal    = trace.SpanKindInternal
	SpanKindServer      = trace.SpanKindServer
	SpanKindClient      = trace.SpanKindClient
	SpanKindProducer    = trace.SpanKindProducer
	SpanKindConsumer    = trace.SpanKindConsumer
)

// otel/trace interface
type Span = trace.Span
type SpanContext = trace.SpanContext

// otel/trace  funcs
var ContextWithRemoteSpanContext = trace.ContextWithRemoteSpanContext
var ContextWithSpan = trace.ContextWithSpan
var ContextWithSpanContext = trace.ContextWithSpanContext
var LinkFromContext = trace.LinkFromContext
var WithSpanKind = trace.WithSpanKind
var WithAttributes = trace.WithAttributes
var SpanFromContext = trace.SpanFromContext
var SpanContextFromContext = trace.SpanContextFromContext

// otel包 常量、类型、方法 快捷链接

// otel types
type ErrorHandler = otel.ErrorHandler
type ErrorHandlerFunc = otel.ErrorHandlerFunc

// 注意此函数改名
var GetTracer = otel.Tracer

// otel funcs
var GetTracerProvider = otel.GetTracerProvider
var SetTracerProvider = otel.SetTracerProvider
var GetTextMapPropagator = otel.GetTextMapPropagator
var SetTextMapPropagator = otel.SetTextMapPropagator
var Handle = otel.Handle
var GetErrorHandler = otel.GetErrorHandler
var SetErrorHandler = otel.SetErrorHandler
var SetLogger = otel.SetLogger
var OtelVersion = otel.Version

func New(config *Config) Provider {
	var exp sdktrace.SpanExporter
	switch config.ExporterName {
	case EXPORTER_NAME_JAEGER:
		exp = newJaegerExporter(config.ExporterEndpoint)
	case EXPORTER_NAME_OLTP:
		exp = newOtlpExporter(config.ExporterEndpoint)
	default:
		exp = NewStdoutExporter()
	}
	return NewProvider(config, exp)
}

func NewProvider(config *Config, exp sdktrace.SpanExporter) Provider {
	res := resource.NewSchemaless(
		semconv.TelemetrySDKLanguageGo,
		semconv.ServiceNameKey.String(config.ServiceName),
	)
	provider := sdktrace.NewTracerProvider(
		// 设置导出exporter
		sdktrace.WithBatcher(exp),
		// 设置采样器
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(config.SampleRate))),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)
	// 设置全局的propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return provider
}

func WithTextMapPropagator(propagator propagation.TextMapPropagator) {
	otel.SetTextMapPropagator(propagator)
}

func NewFileExporter(w io.Writer) sdktrace.SpanExporter {
	if w == nil {
		w = os.Stdout
	}
	exp, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(w),
	)
	if err != nil {
		panic("failed to initialize file trace exporter")
	}
	return exp
}

func NewStdoutExporter() sdktrace.SpanExporter {
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic("failed to initialize stdout trace exporter")
	}
	return exp
}

func newJaegerExporter(url string) sdktrace.SpanExporter {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		panic("failed to initialize jaeger trace exporter")
	}
	return exp
}

func newOtlpExporter(url string) sdktrace.SpanExporter {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		panic("failed to connect to otlp trace exporter")
	}
	exp, err := otlp.New(ctx, otlp.WithGRPCConn(conn))
	if err != nil {
		panic("failed to initialize otlp trace exporter")
	}
	return exp
}
