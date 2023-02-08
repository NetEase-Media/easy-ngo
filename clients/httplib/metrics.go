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
	"net/http"
	"strconv"
	"time"

	"github.com/NetEase-Media/easy-ngo/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// metric norm
// url | method | 调用次数 | 总时间 | 平均rt | 最大并发 | 最慢调用 | 关闭次数 | 错误次数 | 最后错误
// dimension:
// [url] & [method]
// [调用次数] 使用 httplib_request_total counter 单独计算
// [总时间] 也使用 httplib_request_duration_all counter 单独计算吧, 通过下边的NewHistogram得到也可以，不再用单独的点来打
// [平均rt] 使用 [总时间] / 时间 server 端计算
// [最大并发] 也可以通过 [调用次数]的连续的两个点的差值，进行一个大概的估算，计算均值
// [最慢调用] 上报每一次请求的调用时间，可能会漏掉一些点, 修改之前先比较一下，也可以单独计算，这样依赖sdk增加封装的··能力
// [关闭次数] 目前是和 调用次数一样， 目前可以先去掉
// [错误次数] 使用 httplib_request_error 上报各种状态码的错误次数
// [最后错误] 可以把info放到最后的维度之中，但是会有问题，当服务的异常多的时候，会占用大量的内存，会成为系统的缺陷，建议不先实现

// ms0_10 | ms10_100 | ms100_1000 | s1_10 | s10_n
// 使用 NewHistogram 来记录区段的值

// 各种状态码的错误次数
// 使用 [错误次数] 上报的数据为准

// detail:
const (
	metricRequestTotalName         = "httplib_request_total"
	metricRequestDurationName      = "httplib_request_duration"
	metricRequestDurationRangeName = "httplib_request_duration_range"
	metricRequestDurationAllName   = "httplib_request_duration_all"
	metricRequestErrorName         = "httplib_request_error"
)

var (
	metricRequestTotal         metrics.Counter
	metricRequestDuration      metrics.Gauge
	metricRequestDurationRange metrics.Histogram
	metricRequestDurationAll   metrics.Counter
	metricRequestError         metrics.Counter
)

func (c *HttpClient) initMetrics() {
	if c.metrics != nil {
		metricRequestTotal = c.metrics.NewCounter(metricRequestTotalName, "url", "method", "code")
		metricRequestDuration = c.metrics.NewGauge(metricRequestDurationName, "url", "method", "code")
		bukets := prometheus.ExponentialBuckets(.01, 10, 5)
		metricRequestDurationRange = c.metrics.NewHistogram(metricRequestDurationRangeName, bukets, "url", "method", "code")
		metricRequestDurationAll = c.metrics.NewCounter(metricRequestDurationAllName, "url", "method", "code")
		metricRequestError = c.metrics.NewCounter(metricRequestErrorName, "url", "method", "code")
	}
}

func collectCalls(url, method string, start time.Time, code int) {
	du := time.Now().Nanosecond() - start.Nanosecond()/1e6
	fdu := float64(du)
	metricRequestTotal.With("url", url, "method", method, "code", strconv.Itoa(code)).Add(1)
	metricRequestDuration.With("url", url, "method", method, "code", strconv.Itoa(code)).Set(fdu)
	metricRequestDurationRange.With("url", url, "method", method, "code", strconv.Itoa(code)).Observe(fdu)
	metricRequestDurationAll.With("url", url, "method", method, "code", strconv.Itoa(code)).Add(fdu)
}

func collectError(url, method string, code int) {
	if code == http.StatusOK {
		return
	}
	metricRequestError.With("url", url, "method", method, "code", strconv.Itoa(code)).Add(1)
}
