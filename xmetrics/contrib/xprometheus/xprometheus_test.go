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
