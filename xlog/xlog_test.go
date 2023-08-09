package xlog

import "testing"

func TestLogger(t *testing.T) {
	logger := New(DefaultConfig())
	logger.Debugf("debug")
}
