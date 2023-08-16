package app

import (
	"context"
	"sync"
)

var (
	globalPlugins = make(map[Status][]func(ctx context.Context) error)
	mu            = sync.RWMutex{}
)

func RegisterPlugin(status Status, fns ...func(ctx context.Context) error) {
	mu.Lock()
	defer mu.Unlock()
	globalPlugins[status] = append(globalPlugins[status], fns...)
}

func GetFns(status Status) []func(ctx context.Context) error {
	return globalPlugins[status]
}
