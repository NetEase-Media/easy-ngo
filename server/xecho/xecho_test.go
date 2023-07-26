package xecho

import "testing"

func TestServerFromDefault(t *testing.T) {
	s := New(DefaultConfig())
	s.Serve()
}

func TestServerFromEnv(t *testing.T) {
	s := New(EnvConfig())
	s.Serve()
}

func TestServerFromFile(t *testing.T) {
	s := New(FileConfig())
	s.Serve()
}
