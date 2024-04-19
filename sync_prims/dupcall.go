//go:build !solution

package dupcall

import (
	"context"
	"sync"
)

type Call struct {
	ret     interface{}
	err     error
	running chan bool
	ready   chan bool
	close   chan bool
	cnt     int
	mu      sync.Mutex
	once    sync.Once
}

func (o *Call) init() {
	o.running = make(chan bool, 1)
	o.ready = make(chan bool, 1)
	o.close = make(chan bool, 1)
	o.running <- true
}

func (o *Call) run(cb func(context.Context) (interface{}, error)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tmp := make(chan bool, 1)
	go func() {
		o.ret, o.err = cb(ctx)
		tmp <- true
	}()

	select {
	case <-o.close:
		cancel()
	case <-tmp:

	}
	close(o.ready)
	o.ready = make(chan bool, 1)
	o.cnt = 0
	o.running <- true

}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	o.once.Do(o.init)

	select {
	case <-o.running:
		go o.run(cb)
	default:
	}

	o.mu.Lock()
	o.cnt++
	o.mu.Unlock()

	select {
	case <-ctx.Done():
		o.mu.Lock()
		o.cnt--
		o.mu.Unlock()

		// MAYBE UNDER MUTEX
		if o.cnt == 0 {
			o.close <- true
		}

		return nil, ctx.Err()
	case <-o.ready:
		return o.ret, o.err
	}
}
