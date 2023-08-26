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

package xmetrics

type Counter interface {
	With(labelValues ...string) Counter
	Add(delta float64)
	Inc()
}

type Gauge interface {
	With(labelValues ...string) Gauge
	Set(value float64)
	Add(delta float64)
	Inc()
}

type Histogram interface {
	With(labelValues ...string) Histogram
	Observe(value float64)
}

type Provider interface {
	NewCounter(name string, labelNames ...string) Counter
	NewGauge(name string, labelNames ...string) Gauge
	NewHistogram(name string, bucket []float64, labelNames ...string) Histogram
}

type Server interface {
	Stop() error
	Start() error
}

type Bucket struct {
	Start, Factor float64
	Count         int
}
