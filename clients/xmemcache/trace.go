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

package xmemcache

// import (
// 	"context"

// 	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/codes"
// 	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
// 	"go.opentelemetry.io/otel/trace"
// )

// type tracingHook struct {
// }

// func NewTracingHook() *tracingHook {
// 	return &tracingHook{}
// }

// func (h *tracingHook) Before(ctx context.Context, operation string, args ...string) (context.Context, error) {
// 	tr := tracer.GetTracer("memcache")
// 	newCtx, sp := tr.Start(ctx, operation, tracer.WithSpanKind(trace.SpanKindClient))
// 	sp.SetAttributes(
// 		semconv.DBSystemKey.String("memcache"),
// 		semconv.DBOperationKey.String(operation),
// 		attribute.String("db.operation_first_args", args[0]),
// 	)
// 	return newCtx, nil

// }

// func (h *tracingHook) After(ctx context.Context, err error) {
// 	span := tracer.SpanFromContext(ctx)
// 	if span == nil {
// 		return
// 	}
// 	if err != nil {
// 		span.SetStatus(codes.Error, err.Error())
// 	}
// 	span.End()
// }
