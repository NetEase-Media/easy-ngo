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

package httplib

import (
	"context"
	"reflect"

	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
	"github.com/valyala/fasthttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type FastHeaderCarrier struct {
	fasthttp.RequestHeader
	kvs  map[string]string
	keys []string
}

func NewFastHeaderCarrier(rh fasthttp.RequestHeader) *FastHeaderCarrier {
	fhc := &FastHeaderCarrier{
		RequestHeader: rh,
		kvs:           make(map[string]string, rh.Len()),
		keys:          make([]string, rh.Len()),
	}
	rh.VisitAll(func(key, value []byte) {
		fhc.kvs[string(key)] = string(value)
		fhc.keys = append(fhc.keys, string(key))
	})
	return fhc
}

func (fhc FastHeaderCarrier) Get(key string) string {
	return fhc.kvs[key]
}

func (fhc FastHeaderCarrier) Set(key string, value string) {
	fhc.RequestHeader.Set(key, value)
}

func (fhc FastHeaderCarrier) Keys() []string {
	return fhc.keys
}

type tracingHook struct{}

func newTracingHook() {

}

func tracingDoFuc(f DoFunc) DoFunc {
	return func(df *DataFlow, ctx context.Context) (int, error) {
		tr := tracer.GetTracer("httplib")

		propagator := tracer.GetTextMapPropagator()
		//spc := tracer.SpanContextFromContext(ctx)
		//df.logger.Infof("issampled:%s,isvalid:%s,isremote:%s", spc.IsSampled(), spc.IsValid(), spc.IsRemote())
		newCtx, span := tr.Start(ctx, "httpclient", tracer.WithSpanKind(tracer.SpanKindClient))
		df.logger.Debugf("newCtx:%+v,span type:%s, spanCtx:%s", newCtx, reflect.TypeOf(span), span.SpanContext().IsSampled())
		span.SetAttributes(
			semconv.HTTPURLKey.String(string(df.req.URI().FullURI())),
			semconv.HTTPMethodKey.String(string(df.req.Header.Method())),
			semconv.HTTPHostKey.String(string(df.req.Host())),
		)
		propagator.Inject(newCtx, NewFastHeaderCarrier(df.req.Header))

		statusCode, err := f(df, ctx)

		if err != nil {
			span.RecordError(err)
		}

		df.logger.Debugf("end internal,statusCode:%d, spanId:%s, traceId:%s",
			statusCode, span.SpanContext().SpanID(), span.SpanContext().TraceID())
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(statusCode)
		spanCode, spanMsg := semconv.SpanStatusFromHTTPStatusCode(statusCode)
		span.SetAttributes(attrs...)
		span.SetStatus(spanCode, spanMsg)
		span.End()
		return statusCode, err
	}
}
