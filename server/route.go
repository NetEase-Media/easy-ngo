package server

type Route struct {
	Method       METHOD
	RelativePath string
	Handler      any
}
