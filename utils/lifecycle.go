package utils

import (
	"errors"
	"sync"
	"sync/atomic"
)

// Cycle ..
type Cycle struct {
	mu      *sync.Mutex
	wg      *sync.WaitGroup
	done    chan struct{}
	quit    chan error
	closing uint32
	waiting uint32
	// works []func() error
}

// NewCycle new a cycle life
func NewCycle() *Cycle {
	return &Cycle{
		mu:      &sync.Mutex{},
		wg:      &sync.WaitGroup{},
		done:    make(chan struct{}),
		quit:    make(chan error),
		closing: 0,
		waiting: 0,
	}
}

// Run a new goroutine
func (c *Cycle) Run(fn func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.wg.Add(1)
	GoWithRecover(func() {
		defer c.wg.Done()
		if err := fn(); err != nil {
			c.quit <- err
		}
	}, func(r interface{}) {
		c.quit <- errors.New("panic")
	})
}

// Done block and return a chan error
func (c *Cycle) Done() <-chan struct{} {
	if atomic.CompareAndSwapUint32(&c.waiting, 0, 1) {
		go func(c *Cycle) {
			c.mu.Lock()
			defer c.mu.Unlock()
			c.wg.Wait()
			close(c.done)
		}(c)
	}
	return c.done
}

// DoneAndClose ..
func (c *Cycle) DoneAndClose() {
	<-c.Done()
	c.Close()
}

// Close ..
func (c *Cycle) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if atomic.CompareAndSwapUint32(&c.closing, 0, 1) {
		close(c.quit)
	}
}

// Wait blocked for a life cycle
func (c *Cycle) Wait() <-chan error {
	return c.quit
}

func GoWithRecover(handler func(), recoverHandler func(r interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// log.Error(r)
				if recoverHandler != nil {
					go func() {
						defer func() {
							if p := recover(); p != nil {
								// log.Error(p)
							}
						}()
						recoverHandler(r)
					}()
				}
			}
		}()
		handler()
	}()
}
