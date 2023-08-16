package gopool

import (
	"context"
	"fmt"
	"math"
	"sync"
)

// defaultPool is the global default pool.
var defaultPool Pool

var poolMap sync.Map

func init() {
	defaultPool = NewPool("gopool.DefaultPool", math.MaxInt32, NewConfig())
}

// Go is an alternative to the go keyword, which is able to recover panic.
//
//	gopool.Go(func(arg interface{}){
//	    ...
//	}(nil))
func Go(f func()) {
	CtxGo(context.Background(), f)
}

// CtxGo is preferred than Go.
func CtxGo(ctx context.Context, f func()) {
	defaultPool.CtxGo(ctx, f)
}

// SetCap is not recommended to be called, this func changes the global pool's capacity which will affect other callers.
func SetCap(cap int32) {
	defaultPool.SetCap(cap)
}

// SetPanicHandler sets the panic handler for the global pool.
func SetPanicHandler(f func(context.Context, interface{})) {
	defaultPool.SetPanicHandler(f)
}

// WorkerCount returns the number of global default pool's running workers
func WorkerCount() int32 {
	return defaultPool.WorkerCount()
}

// RegisterPool registers a new pool to the global map.
// GetPool can be used to get the registered pool by name.
// returns error if the same name is registered.
func RegisterPool(p Pool) error {
	_, loaded := poolMap.LoadOrStore(p.Name(), p)
	if loaded {
		return fmt.Errorf("name: %s already registered", p.Name())
	}
	return nil
}

// GetPool gets the registered pool by name.
// Returns nil if not registered.
func GetPool(name string) Pool {
	p, ok := poolMap.Load(name)
	if !ok {
		return nil
	}
	return p.(Pool)
}
