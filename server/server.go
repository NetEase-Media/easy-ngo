package server

import "net/http"

const (
	GET     = http.MethodGet
	HEAD    = http.MethodHead
	POST    = http.MethodPost
	PUT     = http.MethodPut
	PATCH   = http.MethodPatch
	DELETE  = http.MethodDelete
	CONNECT = http.MethodConnect
	OPTIONS = http.MethodOptions
	TRACE   = http.MethodTrace
)

type Server interface {
	Serve() error
	Shutdown() error
	Healthz() bool
	Init() error
}
