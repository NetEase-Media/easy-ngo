package main

import (
	"fmt"
	"os"

	"github.com/NetEase-Media/easy-ngo/config"
	_ "github.com/NetEase-Media/easy-ngo/config/contrib/env"
)

func main() {
	Config()
}

func Config() {
	os.Setenv("APP_NAME", "easy-ngo")
	os.Setenv("TEST_NAME1", "test-name1")
	os.Setenv("TEST_NAME2", "test-name2")
	c := config.New()
	c.AddProtocol("env://prefix=APP;envName=TEST_NAME1,TEST_NAME2")
	c.AddProtocol("file://name=config;type=yaml;path=./config.yaml,./config2.yaml")
	c.Init()
	err := c.ReadConfig()
	if err != nil {
		fmt.Print(err)
	}
	config.WithConfig(c)
	fmt.Println("APP_NAME=" + config.GetString("name"))
	fmt.Println("TEST_NAME1=" + config.GetString("test_name1"))
	fmt.Println("TEST_NAME2=" + config.GetString("test_name2"))
}
