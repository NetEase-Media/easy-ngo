package server

type Server interface {
	Serve() error
	Shutdown() error
	Healthz() bool
}
