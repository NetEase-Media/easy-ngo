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
	"net/http"

	"github.com/NetEase-Media/easy-ngo/xmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusProvider struct {
	config *Config
}

func NewProvider(config *Config) xmetrics.Provider {
	return &prometheusProvider{
		config: config,
	}
}

func (p *prometheusProvider) NewCounter(name string, labelNames ...string) xmetrics.Counter {
	return NewCounterFrom(prometheus.CounterOpts{
		Namespace: p.config.Namespace,
		Subsystem: p.config.Subsystem,
		Name:      name,
		Help:      name,
	}, labelNames)
}

func (p *prometheusProvider) NewGauge(name string, labelNames ...string) xmetrics.Gauge {
	return NewGaugeFrom(prometheus.GaugeOpts{
		Namespace: p.config.Namespace,
		Subsystem: p.config.Subsystem,
		Name:      name,
		Help:      name,
	}, labelNames)
}

func (p *prometheusProvider) NewHistogram(name string, buckets []float64, labelNames ...string) xmetrics.Histogram {
	return NewHistogramFrom(prometheus.HistogramOpts{
		Namespace: p.config.Namespace,
		Subsystem: p.config.Subsystem,
		Name:      name,
		Help:      name,
		Buckets:   buckets,
	}, labelNames)
}

func (p *prometheusProvider) Stop() {}

func (p *prometheusProvider) Start() error {
	http.Handle(p.config.Path, promhttp.Handler())
	return http.ListenAndServe(p.config.Addr, nil)
}
