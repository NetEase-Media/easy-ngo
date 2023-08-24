package pluginxmetrics

import "github.com/NetEase-Media/easy-ngo/xmetrics"

var provider xmetrics.Provider

func NewCounter(name string, labelNames ...string) xmetrics.Counter {
	return provider.NewCounter(name, labelNames...)
}

func NewGauge(name string, labelNames ...string) xmetrics.Gauge {
	return provider.NewGauge(name, labelNames...)
}

func NewHistogram(name string, buckets []float64, labelNames ...string) xmetrics.Histogram {
	return provider.NewHistogram(name, buckets, labelNames...)
}

func WithVendor(p xmetrics.Provider) {
	provider = p
}

func GetProvider() xmetrics.Provider {
	return provider
}
