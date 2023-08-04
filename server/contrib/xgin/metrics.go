package xgin

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/NetEase-Media/easy-ngo/xmetrics"
)

var (
	metricRequestTotal    = "gin_request_total"
	metricRequestUVTotal  = "gin_request_uv_total"
	metricURIRequestTotal = "gin_uri_request_total"
	metricRequestBody     = "gin_request_body_total"
	metricResponseBody    = "gin_response_body_total"
	metricRequestDuration = "gin_request_duration"
	metricSlowRequest     = "gin_slow_request_total"
)

var (
	requestTotal    xmetrics.Counter
	requestDuration xmetrics.Histogram
)

func (server *Server) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if server.config.MetricsPath == c.Request.URL.Path {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		calc(start, c)
	}
}

func (server *Server) initMetrics() {
	// bukets := prometheus.ExponentialBuckets(10, 10, 5)
	// requestTotal = xmetrics.NewCounter(metricRequestTotal, "url", "method", "code")
	// requestDuration = xmetrics.NewHistogram(metricRequestDuration, bukets, "url", "method", "code")
}

func calc(start time.Time, c *gin.Context) {
	requestTotal.With("url", c.Request.URL.Path, "method", c.Request.Method, "code", strconv.Itoa(c.Writer.Status())).Add(1)
	duration := (time.Now().Nanosecond() - start.Nanosecond()) / 1e6
	requestDuration.With("url", c.Request.URL.Path, "method", c.Request.Method, "code", strconv.Itoa(c.Writer.Status())).Observe(float64(duration))
}
