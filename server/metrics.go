package server

import (
	"strconv"

	"github.com/NetEase-Media/easy-ngo/xmetrics"
)

var (
	requestTotal xmetrics.Counter
)

var (
	metricRequestTotal    = "gin_request_total"
	metricRequestUVTotal  = "gin_request_uv_total"
	metricURIRequestTotal = "gin_uri_request_total"
	metricRequestBody     = "gin_request_body_total"
	metricResponseBody    = "gin_response_body_total"
	metricRequestDuration = "gin_request_duration"
	metricSlowRequest     = "gin_slow_request_total"

	labelDomain = "domain"
	labelURL    = "url"
	labelMethod = "method"
	labelCode   = "code"
)

type HttpMetrics struct {
}

type HttpLabels struct {
	Url    string
	Method string
	Code   int
	Domain string
}

func (httpMetrics *HttpMetrics) Record(duration int, labels HttpLabels) {
	requestTotal.With(labelDomain, labels.Domain, labelURL, labels.Url, labelMethod, labels.Method, labelCode, strconv.Itoa(labels.Code)).Add(1)
}

func NewHttpMetrics() *HttpMetrics {
	return &HttpMetrics{}
}

func (httpMetrics *HttpMetrics) Init() {
	// bukets := prometheus.ExponentialBuckets(10, 10, 5)
	// requestTotal = xmetrics.NewCounter(metricRequestTotal, "url", "method", "code")
	// requestDuration = xmetrics.NewHistogram(metricRequestDuration, bukets, "url", "method", "code")
}
