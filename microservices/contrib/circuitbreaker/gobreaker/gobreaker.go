package gobreaker

import (
	"errors"

	"github.com/NetEase-Media/easy-ngo/microservices/middleware/circuitbreak"
	"github.com/sony/gobreaker"
)

var _ circuitbreak.CircuitBreaker = (*CircuitBreaker)(nil)

func NewCircuitBreaker(sts gobreaker.Settings) circuitbreak.CircuitBreaker {
	return &CircuitBreaker{
		CB: gobreaker.NewCircuitBreaker(sts),
	}
}

// CircuitBreaker is a circuit breaker.
type CircuitBreaker struct {
	CB *gobreaker.CircuitBreaker
}

func (c *CircuitBreaker) Execute(f func() error) bool {
	_, err := c.CB.Execute(func() (interface{}, error) {
		return nil, f()
	})
	if err != nil && (errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests)) {
		return true
	}
	return false
}
