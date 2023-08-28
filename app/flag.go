package app

import (
	"flag"
	"strings"
)

type arrayFlags []string

var configProtocols arrayFlags

func (s *arrayFlags) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s *arrayFlags) String() string {
	return strings.Join(*s, ",")
}

func parse() {
	flag.Var(&configProtocols, "c", "Config file list!")
	flag.Parse()
}

func GetConfigProtocols() []string {
	return configProtocols
}
