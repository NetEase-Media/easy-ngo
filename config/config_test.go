package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	// _ "github.com/NetEase-Media/easy-ngo/config/contrib/env"
)

func TestEnv(t *testing.T) {
	os.Setenv("APP_NAME", "easy-ngo")
	New().WithConfigArgument("env://prefix=APP")
	config.Init()
	err := config.ReadConfig()
	if err != nil {
		fmt.Print(err)
	}
	v := GetString("name")
	assert.Equal(t, "easy-ngo", v)
}
