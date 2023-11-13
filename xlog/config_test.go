package xlog

import (
	"testing"

	"gotest.tools/assert"
)

func TestConfig(t *testing.T) {
	c := DefaultConfig()
	assert.Equal(t, c.Level, INFO, "info error")
	assert.Equal(t, c.Format, JSON, "format error")
	assert.Equal(t, c.Path, "./logs", "path error")
	assert.Equal(t, c.ErrorPath, "./logs", "error path error")
	assert.Equal(t, c.FileName, "app", "filename error")
	assert.Equal(t, c.MaxAge, 7, "max age error")
	assert.Equal(t, c.MaxBackups, 7, "max backups error")
	assert.Equal(t, c.MaxSize, 1024, "max size error")
	assert.Equal(t, c.Console, true, "console error")
	assert.Equal(t, c.Suffix, "log", "suffix error")
	assert.Equal(t, c.ErrorSuffix, "error.log", "error suffix error")
	assert.Equal(t, c.ErrlogLevel, ERROR, "error log level error")
}
