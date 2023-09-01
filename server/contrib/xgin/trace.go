package xgin

import (
	"strings"

	"github.com/NetEase-Media/easy-ngo/server"
	"github.com/NetEase-Media/easy-ngo/xtracer"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.16.0"
)

var (
	gtracer xtracer.Tracer
)

func (s *Server) traceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/health") {
			propagator := xtracer.GetTextMapPropagator()
			oldContext := c.Request.Context()
			ctx := propagator.Extract(oldContext, propagation.HeaderCarrier(c.Request.Header))
			ctx, span := gtracer.Start(
				ctx, c.FullPath(),
				xtracer.WithSpanKind(xtracer.SpanKindServer),
			)
			span.SetAttributes(semconv.HTTPURLKey.String(c.FullPath()))
			span.SetAttributes(semconv.HostNameKey.String(c.Request.Host))
			span.SetAttributes(semconv.HTTPMethodKey.String(c.Request.Method))
			span.SetAttributes(semconv.HTTPRequestContentLengthKey.Int64(c.Request.ContentLength))
			// 替换 context
			c.Request = c.Request.WithContext(ctx)
			defer func() {
				// 记录response
				code := c.Writer.Status()
				spanCode, spanMsg := server.SpanStatusFromHTTPStatusCode(code)
				span.SetAttributes(semconv.HTTPResponseContentLengthKey.String(c.Writer.Header().Get("Content-Length")))
				span.SetStatus(spanCode, spanMsg)
				span.End()
			}()
		}
		c.Next()
	}
}

func (server *Server) initTracer() {
	gtracer = xtracer.GetTracer("gin")
}
