package filter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddInclude(t *testing.T) {
	var (
		path  string = "/test"
		path2 string = "/test2"
	)

	f := NewDefaultFilter(&Options{
		Include: []string{path},
	})

	err := f.AddInclude(context.Background(), path2)
	if err != nil {
		t.Errorf("AddInclude() error = %v", err)
	}
}

func TestExclude(t *testing.T) {
	var (
		path  string = "/test"
		path2 string = "/test2"
	)

	f := NewDefaultFilter(&Options{
		Exclude: []string{path},
	})

	err := f.AddExclude(context.Background(), path2)
	if err != nil {
		t.Errorf("AddExclude() error = %v", err)
	}
}

func TestFind(t *testing.T) {
	var (
		path  string = "/test"
		path2 string = `/a/:name`
	)

	f := NewDefaultFilter(&Options{
		Include: []string{path, path2},
	})

	ok, _ := f.DoFilt(context.Background(), `/a/bc`)
	assert.Equal(t, false, ok) // 因为在include中，所以不需要被过滤掉
	ok, _ = f.DoFilt(context.Background(), `/a/bc/d`)
	assert.Equal(t, false, ok) // 因为不在include中， 也不在exclude当中，所以不需要被过滤掉
}
