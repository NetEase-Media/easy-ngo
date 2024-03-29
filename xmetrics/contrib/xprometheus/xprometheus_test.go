package xprometheus

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"gotest.tools/assert"
)

func TestCounter(t *testing.T) {
	// with label
	c := NewCounterFrom(prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "http_requests_total",
		Help:      "http_requests_total",
	}, []string{"method", "code"})

	c.With("method", "GET", "code", "200").Add(2)
	value := make(chan prometheus.Metric, 1)
	c.cv.Collect(value)
	m := <-value
	data := dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 2.0, *data.Counter.Value)

	c.With("method", "GET", "code", "200").Inc()
	value = make(chan prometheus.Metric, 1)
	c.cv.Collect(value)
	m = <-value
	data = dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 3.0, *data.Counter.Value)

	// without label
	c2 := NewCounterFrom(prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "redis_requests_total",
		Help:      "redis_requests_total",
	}, nil)

	c2.Add(2)
	value = make(chan prometheus.Metric, 1)
	c2.cv.Collect(value)
	m = <-value
	data = dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 2.0, *data.Counter.Value)
}

func TestGauge(t *testing.T) {
	// with label
	g := NewGaugeFrom(prometheus.GaugeOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "cup_use",
		Help:      "cup_use",
	}, []string{"host"})

	g.With("host", "127.0.0.1").Set(2)
	value := make(chan prometheus.Metric, 1)
	g.gv.Collect(value)
	m := <-value
	data := dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 2.0, *data.Gauge.Value)

	g.With("host", "127.0.0.1").Inc()
	value = make(chan prometheus.Metric, 1)
	g.gv.Collect(value)
	m = <-value
	data = dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 3.0, *data.Gauge.Value)

	g.With("host", "127.0.0.1").Add(2)
	value = make(chan prometheus.Metric, 1)
	g.gv.Collect(value)
	m = <-value
	data = dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 5.0, *data.Gauge.Value)
}

func TestSummary(t *testing.T) {
	// with label
	s := NewSummaryFrom(prometheus.SummaryOpts{
		Namespace:  "",
		Subsystem:  "",
		Name:       "http_requests_duration",
		Help:       "http_requests_duration",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"method", "code"})

	s.With("method", "GET", "code", "200").Observe(2)
	value := make(chan prometheus.Metric, 1)
	s.sv.Collect(value)
	m := <-value
	data := dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 2.0, *data.Summary.SampleSum)
	assert.Equal(t, uint64(1), *data.Summary.SampleCount)

	s.With("method", "GET", "code", "200").Observe(3)
	value = make(chan prometheus.Metric, 1)
	s.sv.Collect(value)
	m = <-value
	data = dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 5.0, *data.Summary.SampleSum)
	assert.Equal(t, uint64(2), *data.Summary.SampleCount)

	// without label
	s2 := NewSummaryFrom(prometheus.SummaryOpts{
		Namespace:  "",
		Subsystem:  "",
		Name:       "redis_requests_total",
		Help:       "redis_requests_total",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, nil)

	s2.Observe(2)
	value = make(chan prometheus.Metric, 1)
	s2.sv.Collect(value)
	m = <-value
	data = dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 2.0, *data.Summary.SampleSum)
	assert.Equal(t, uint64(1), *data.Summary.SampleCount)
}

func TestHistogram(t *testing.T) {
	// with label
	h := NewHistogramFrom(prometheus.HistogramOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "http_requests_duration",
		Help:      "http_requests_duration",
		Buckets:   []float64{10, 200, 400, 600, 1000},
	}, []string{"method", "code"})

	h.With("method", "GET", "code", "200").Observe(100)
	value := make(chan prometheus.Metric, 1)
	h.hv.Collect(value)
	m := <-value
	data := dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 100.0, *data.Histogram.SampleSum)
	assert.Equal(t, uint64(1), *data.Histogram.SampleCount)

	h.With("method", "GET", "code", "200").Observe(200)
	value = make(chan prometheus.Metric, 1)
	h.hv.Collect(value)
	m = <-value
	data = dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 300.0, *data.Histogram.SampleSum)
	assert.Equal(t, uint64(2), *data.Histogram.SampleCount)

	// without label
	h2 := NewHistogramFrom(prometheus.HistogramOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "redis_requests_total",
		Help:      "redis_requests_total",
		Buckets:   []float64{10, 200, 400, 600, 1000},
	}, nil)

	h2.Observe(100)
	value = make(chan prometheus.Metric, 1)
	h2.hv.Collect(value)
	m = <-value
	data = dto.Metric{}
	m.Write(&data)
	assert.Equal(t, 100.0, *data.Histogram.SampleSum)
	assert.Equal(t, uint64(1), *data.Histogram.SampleCount)
}
