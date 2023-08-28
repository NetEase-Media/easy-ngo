package server

import (
	"net/http"
)

type METHOD string

const (
	GET     METHOD = http.MethodGet
	HEAD           = http.MethodHead
	POST           = http.MethodPost
	PUT            = http.MethodPut
	PATCH          = http.MethodPatch
	DELETE         = http.MethodDelete
	CONNECT        = http.MethodConnect
	OPTIONS        = http.MethodOptions
	TRACE          = http.MethodTrace
)

type Server interface {
	Serve() error
	Shutdown() error
	Healthz() bool
	Init() error

	// GET(relativePath string, handler any)
	// POST(relativePath string, handler any)
	// PUT(relativePath string, handler any)
	// DELETE(relativePath string, handler any)
	// PATCH(relativePath string, handler any)
	// HEAD(relativePath string, handler any)
	// OPTIONS(relativePath string, handler any)
}
