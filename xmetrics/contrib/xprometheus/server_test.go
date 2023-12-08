package xprometheus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	s := NewServer(DefaultConfig())
	assert.NotNil(t, s)
}
