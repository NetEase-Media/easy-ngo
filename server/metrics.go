package server

import (
	"strconv"
	"time"

	"github.com/NetEase-Media/easy-ngo/xmetrics"
)

var (
	requestTotal    xmetrics.Counter
	requestDuration xmetrics.Histogram
)

var (
	metricRequestTotal    = "request_total"
	metricRequestDuration = "request_duration"

	LABELDOMAIN = "domain"
	LABELURL    = "url"
	LABELMETHOD = "method"
	LABELCODE   = "code"
)

type HttpMetrics struct {
	metrics xmetrics.Provider
	bucket  xmetrics.Bucket
}

type HttpLabels struct {
	Url    string
	Method string
	Code   int
	Domain string
}

func (httpMetrics *HttpMetrics) Record(labels HttpLabels, start time.Time, end time.Time) {
	requestTotal.With(LABELDOMAIN, labels.Domain, LABELURL, labels.Url, LABELMETHOD, labels.Method, LABELCODE, strconv.Itoa(labels.Code)).Add(1)
	requestDuration.With(LABELDOMAIN, labels.Domain, LABELURL, labels.Url, LABELMETHOD, labels.Method, LABELCODE).Observe(float64((end.Nanosecond() - start.Nanosecond()) / 1e6))
}

func NewHttpMetrics(metrics xmetrics.Provider, bucket xmetrics.Bucket) *HttpMetrics {
	return &HttpMetrics{
		metrics: metrics,
		bucket:  bucket,
	}
}

func (httpMetrics *HttpMetrics) Init() {
	requestTotal = httpMetrics.metrics.NewCounter(metricRequestTotal, LABELDOMAIN, LABELURL, LABELMETHOD, LABELCODE)
	requestDuration = httpMetrics.metrics.NewHistogram(metricRequestDuration, httpMetrics.exponentialBuckets(httpMetrics.bucket.Start, httpMetrics.bucket.Factor, httpMetrics.bucket.Count), LABELDOMAIN, LABELURL, LABELMETHOD, LABELCODE)
}

func (httpMetrics *HttpMetrics) exponentialBuckets(start, factor float64, count int) []float64 {
	if count < 1 {
		panic("ExponentialBuckets needs a positive count")
	}
	if start <= 0 {
		panic("ExponentialBuckets needs a positive start value")
	}
	if factor <= 1 {
		panic("ExponentialBuckets needs a factor greater than 1")
	}
	buckets := make([]float64, count)
	for i := range buckets {
		buckets[i] = start
		start *= factor
	}
	return buckets
}
