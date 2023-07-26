package server

type Server interface {
	Serve() error
	Shutdown() error
	Healthz() bool
	Info() *Info
}

type Info struct {
}

type Config struct {
	Host           string
	Port           int
	Mode           string
	EnabledMetric  bool
	EnableTracer   bool
	ServiceAddress string
	MetricsPath    string
}

type Healthz struct {
}
