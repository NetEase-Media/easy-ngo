package filter

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddIncludeRoute(t *testing.T) {
	var (
		path  string = "/test"
		path2 string = "/test2"
	)

	f := NewDefaultRouteFilter(&RouteOptions{
		Include: []string{path},
	})

	err := f.AddInclude(context.Background(), path2)
	if err != nil {
		t.Errorf("AddInclude() error = %v", err)
	}
}

func TestExcludeRoute(t *testing.T) {
	var (
		path  string = "/test"
		path2 string = "/test2"
	)

	f := NewDefaultRouteFilter(&RouteOptions{
		Exclude: []string{path},
	})

	err := f.AddExclude(context.Background(), path2)
	if err != nil {
		t.Errorf("AddExclude() error = %v", err)
	}
}

func TestFindRoute(t *testing.T) {
	var (
		path  string = "/test"
		path2 string = "/eat/{v1:[a-zA-Z]{3}[0-1]{3}}"
	)

	f := NewDefaultRouteFilter(&RouteOptions{
		Include: []string{path, path2},
	})
	req, _ := http.NewRequest("GET", "/eat/aaa/888", nil)
	match := new(mux.RouteMatch)
	ok := f.BelongToInclude(req, match)
	assert.Equal(t, false, ok)
	ok = f.BelongToExclude(req, match)
	assert.Equal(t, false, ok)

	req2, _ := http.NewRequest("GET", "/eat/aaa111", nil)
	ok = f.BelongToInclude(req2, match)
	assert.Equal(t, true, ok)
	ok = f.BelongToExclude(req2, match)
	assert.Equal(t, false, ok)
}
