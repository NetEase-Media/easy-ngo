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
	"strconv"
	"time"

	"github.com/NetEase-Media/easy-ngo/observability/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
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
	requestTotal    metrics.Counter
	requestDuration metrics.Histogram
)

func (server *Server) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if server.option.MetricsPath == c.Request.URL.Path {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		calc(start, c)
	}
}

func (server *Server) initMetrics() {
	bukets := prometheus.ExponentialBuckets(10, 10, 5)
	requestTotal = server.Metrics.NewCounter(metricRequestTotal, "url", "method", "code")
	requestDuration = server.Metrics.NewHistogram(metricRequestDuration, bukets, "url", "method", "code")
}

func calc(start time.Time, c *gin.Context) {
	requestTotal.With("url", c.Request.URL.Path, "method", c.Request.Method, "code", strconv.Itoa(c.Writer.Status())).Add(1)
	duration := (time.Now().Nanosecond() - start.Nanosecond()) / 1e6
	requestDuration.With("url", c.Request.URL.Path, "method", c.Request.Method, "code", strconv.Itoa(c.Writer.Status())).Observe(float64(duration))
}
