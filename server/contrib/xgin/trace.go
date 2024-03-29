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

package xgin

import (
	"strings"

	"github.com/NetEase-Media/easy-ngo/xtracer"
	"github.com/gin-gonic/gin"
)

var (
	gtracer xtracer.Tracer
)

func (server *Server) traceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/health") {
			// propagator := xtracer.GetTextMapPropagator()
			// oldContext := c.Request.Context()
			// ctx := propagator.Extract(oldContext, propagation.HeaderCarrier(c.Request.Header))
			// ctx, span := gtracer.Start(
			// 	ctx, c.FullPath(),
			// 	xtracer.WithSpanKind(xtracer.SpanKindServer),
			// )
			// kvs := semconv.NetAttributesFromHTTPRequest("tcp", c.Request)
			// kvs = append(kvs, semconv.HTTPURLKey.String(c.FullPath()))
			// kvs = append(kvs, semconv.HTTPHostKey.String(c.Request.Host))
			// kvs = append(kvs, semconv.HTTPMethodKey.String(c.Request.Method))
			// kvs = append(kvs, semconv.HTTPRequestContentLengthKey.Int64(c.Request.ContentLength))
			// span.SetAttributes(kvs...)
			// // 替换 context
			// c.Request = c.Request.WithContext(ctx)
			// defer func() {
			// 	// 记录response
			// 	code := c.Writer.Status()
			// 	attrs := semconv.HTTPAttributesFromHTTPStatusCode(code)
			// 	spanCode, spanMsg := semconv.SpanStatusFromHTTPStatusCode(code)
			// 	span.SetAttributes(attrs...)
			// 	span.SetStatus(spanCode, spanMsg)
			// 	span.End()
			// }()
		}
		c.Next()
	}
}

func (server *Server) initTracer() {
	gtracer = xtracer.GetTracer("gin")
}
