package server

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/codes"
)

type METHOD string

type codeRange struct {
	fromInclusive int
	toInclusive   int
}

var validRangesPerCategory = map[int][]codeRange{
	1: {
		{http.StatusContinue, http.StatusEarlyHints},
	},
	2: {
		{http.StatusOK, http.StatusAlreadyReported},
		{http.StatusIMUsed, http.StatusIMUsed},
	},
	3: {
		{http.StatusMultipleChoices, http.StatusUseProxy},
		{http.StatusTemporaryRedirect, http.StatusPermanentRedirect},
	},
	4: {
		{http.StatusBadRequest, http.StatusTeapot}, // yes, teapot is so useful…
		{http.StatusMisdirectedRequest, http.StatusUpgradeRequired},
		{http.StatusPreconditionRequired, http.StatusTooManyRequests},
		{http.StatusRequestHeaderFieldsTooLarge, http.StatusRequestHeaderFieldsTooLarge},
		{http.StatusUnavailableForLegalReasons, http.StatusUnavailableForLegalReasons},
	},
	5: {
		{http.StatusInternalServerError, http.StatusLoopDetected},
		{http.StatusNotExtended, http.StatusNetworkAuthenticationRequired},
	},
}

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
}

func SpanStatusFromHTTPStatusCode(code int) (codes.Code, string) {
	spanCode, valid := validateHTTPStatusCode(code)
	if !valid {
		return spanCode, fmt.Sprintf("Invalid HTTP status code %d", code)
	}
	return spanCode, ""
}

func validateHTTPStatusCode(code int) (codes.Code, bool) {
	category := code / 100
	ranges, ok := validRangesPerCategory[category]
	if !ok {
		return codes.Error, false
	}
	ok = false
	for _, crange := range ranges {
		ok = crange.contains(code)
		if ok {
			break
		}
	}
	if !ok {
		return codes.Error, false
	}
	if category > 0 && category < 4 {
		return codes.Unset, true
	}
	return codes.Error, true
}

func (r codeRange) contains(code int) bool {
	return r.fromInclusive <= code && code <= r.toInclusive
}
