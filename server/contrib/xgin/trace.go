package xgin

import (
	"strings"

	tracer "github.com/NetEase-Media/easy-ngo/observability/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

var (
	gtracer tracer.Tracer
)

func (server *Server) traceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/health") {
			propagator := tracer.GetTextMapPropagator()
			oldContext := c.Request.Context()
			ctx := propagator.Extract(oldContext, propagation.HeaderCarrier(c.Request.Header))
			ctx, span := gtracer.Start(
				ctx, c.FullPath(),
				tracer.WithSpanKind(tracer.SpanKindServer),
			)
			kvs := semconv.NetAttributesFromHTTPRequest("tcp", c.Request)
			kvs = append(kvs, semconv.HTTPURLKey.String(c.FullPath()))
			kvs = append(kvs, semconv.HTTPHostKey.String(c.Request.Host))
			kvs = append(kvs, semconv.HTTPMethodKey.String(c.Request.Method))
			kvs = append(kvs, semconv.HTTPRequestContentLengthKey.Int64(c.Request.ContentLength))
			span.SetAttributes(kvs...)
			// 替换 context
			c.Request = c.Request.WithContext(ctx)
			defer func() {
				// 记录response
				code := c.Writer.Status()
				attrs := semconv.HTTPAttributesFromHTTPStatusCode(code)
				spanCode, spanMsg := semconv.SpanStatusFromHTTPStatusCode(code)
				span.SetAttributes(attrs...)
				span.SetStatus(spanCode, spanMsg)
				span.End()
			}()
		}
		c.Next()
	}
}

func (server *Server) initTracer() {
	gtracer = tracer.GetTracer("gin")
}
