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
