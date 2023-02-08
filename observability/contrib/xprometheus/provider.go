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

package xprometheus

import (
	"github.com/NetEase-Media/easy-ngo/observability/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type prometheusProvider struct {
	namespace string
	subsystem string
}

func NewProvider(namespace, subsystem string) metrics.Provider {
	return &prometheusProvider{
		namespace: namespace,
		subsystem: subsystem,
	}
}

func DefaultProvider() metrics.Provider {
	return &prometheusProvider{
		namespace: "",
		subsystem: "",
	}
}

func (p *prometheusProvider) NewCounter(name string, labelNames ...string) metrics.Counter {
	return NewCounterFrom(prometheus.CounterOpts{
		Namespace: p.namespace,
		Subsystem: p.subsystem,
		Name:      name,
		Help:      name,
	}, labelNames)
}

func (p *prometheusProvider) NewGauge(name string, labelNames ...string) metrics.Gauge {
	return NewGaugeFrom(prometheus.GaugeOpts{
		Namespace: p.namespace,
		Subsystem: p.subsystem,
		Name:      name,
		Help:      name,
	}, labelNames)
}

func (p *prometheusProvider) NewHistogram(name string, buckets []float64, labelNames ...string) metrics.Histogram {
	return NewHistogramFrom(prometheus.HistogramOpts{
		Namespace: p.namespace,
		Subsystem: p.subsystem,
		Name:      name,
		Help:      name,
		Buckets:   buckets,
	}, labelNames)
}

func (p *prometheusProvider) Stop() {}
