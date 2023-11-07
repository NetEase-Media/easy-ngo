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
