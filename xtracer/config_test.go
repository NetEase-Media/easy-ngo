package xtracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefa(t *testing.T) {
	c := DefaultConfig()
	assert.Equal(t, EXPORTER_NAME_STDOUT, c.ExporterName)
	assert.Equal(t, 1.0, c.SampleRate)
}
