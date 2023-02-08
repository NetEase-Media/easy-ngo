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

package xredis

import (
	"context"
	"strings"

	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
	"github.com/go-redis/redis/extra/rediscmd"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type tracingHook struct {
	container *RedisContainer
}

func newTracingHook(container *RedisContainer) *tracingHook {
	return &tracingHook{container: container}
}

func (th *tracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	//TODO add enable
	newCtx, span := tracer.GetTracer("redis").Start(
		ctx, "redis:"+cmd.FullName(),
		tracer.WithSpanKind(tracer.SpanKindClient),
	)
	addrs := strings.Join(th.container.Opt.Addr, ",")
	span.SetAttributes(
		semconv.DBSystemRedis,
		semconv.NetHostPortKey.String(addrs),
		semconv.DBOperationKey.String(cmd.Name()),
	)
	return newCtx, nil
}
func (th *tracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	//TODO add enable
	span := tracer.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	if err := cmd.Err(); err != nil && err != redis.Nil {
		recordError(ctx, span, err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "success")
	}
	span.End()
	return nil
}

func (th *tracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	//TODO add enable
	summary, cmdsString := rediscmd.CmdsString(cmds)

	newCtx, span := tracer.GetTracer("redis").Start(
		ctx, "redis-pipeline:"+summary,
		tracer.WithSpanKind(tracer.SpanKindClient),
	)
	addrs := strings.Join(th.container.Opt.Addr, ",")
	span.SetAttributes(
		semconv.DBSystemRedis,
		semconv.NetHostPortKey.String(addrs),
		semconv.DBOperationKey.String(cmdsString),
	)

	return newCtx, nil
}

func (th *tracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	//TODO add enable
	span := tracer.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	span.SetStatus(codes.Ok, "success")
	if err := cmds[0].Err(); err != nil && err != redis.Nil {
		recordError(ctx, span, err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
	return nil
}

func recordError(ctx context.Context, span tracer.Span, err error) {
	if err != redis.Nil {
		span.RecordError(err)
	}
}
