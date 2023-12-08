package xprometheus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProvider(t *testing.T) {
	p := NewProvider(DefaultConfig())
	assert.NotNil(t, p)
}

func TestNewCounter(t *testing.T) {
	p := NewProvider(DefaultConfig())
	assert.NotNil(t, p)

	c := p.NewCounter("http_requests_total", "method", "code")
	assert.NotNil(t, c)

	c2 := p.NewCounter("redis_requests_total", "operation")
	assert.NotNil(t, c2)

	c3 := p.NewCounter("db_request_total")
	assert.NotNil(t, c3)
}

func TestNewGauge(t *testing.T) {
	p := NewProvider(DefaultConfig())
	assert.NotNil(t, p)

	g := p.NewGauge("cup_use", "host")
	assert.NotNil(t, g)

	g2 := p.NewGauge("memery_use_total")
	assert.NotNil(t, g2)
}

func TestNewHistogram(t *testing.T) {
	p := NewProvider(DefaultConfig())
	assert.NotNil(t, p)

	h := p.NewHistogram("http_request_duration_seconds", []float64{0.1, 0.2, 0.3, 0.4, 0.5}, "method", "code")
	assert.NotNil(t, h)

	h2 := p.NewHistogram("redis_request_duration_seconds", []float64{0.1, 0.2, 0.3, 0.4, 0.5}, "operation")
	assert.NotNil(t, h2)

	h3 := p.NewHistogram("db_request_duration_seconds", []float64{0.1, 0.2, 0.3, 0.4, 0.5})
	assert.NotNil(t, h3)
}
