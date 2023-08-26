package xprometheus

import (
	"net/http"

	"github.com/NetEase-Media/easy-ngo/xmetrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheuServer struct {
	config *Config
}

func NewServer(config *Config) xmetrics.Server {
	return &prometheuServer{
		config: config,
	}
}

func (p *prometheuServer) Stop() error {
	return nil
}

func (p *prometheuServer) Start() error {
	http.Handle(p.config.Path, promhttp.Handler())
	return http.ListenAndServe(p.config.Addr, nil)
}
