package xmetrics

var provider Provider

func NewCounter(name string, labelNames ...string) Counter {
	return provider.NewCounter(name, labelNames...)
}

func NewGauge(name string, labelNames ...string) Gauge {
	return provider.NewGauge(name, labelNames...)
}

func NewHistogram(name string, buckets []float64, labelNames ...string) Histogram {
	return provider.NewHistogram(name, buckets, labelNames...)
}

func WithVendor(p Provider) {
	provider = p
}

func GetProvider() Provider {
	return provider
}
