package xprometheus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()

	assert.Equal(t, "", c.Namespace)
	assert.Equal(t, "", c.Subsystem)
	assert.Equal(t, "/metrics", c.Path)
	assert.Equal(t, ":8888", c.Addr)
}
