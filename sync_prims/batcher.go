//go:build !solution

package batcher

import (
	"sync"

	"gitlab.com/slon/shad-go/batcher/slow"
	"golang.org/x/sync/singleflight"
)

type Batcher struct {
	sf   *singleflight.Group
	v    *slow.Value
	mu   *sync.Mutex
	cond *sync.Cond
	once *sync.Once
	wait bool
}

func (b *Batcher) Load() (v interface{}) {

	b.mu.Lock()
	if b.wait {
		b.cond.Wait()
	}
	b.wait = true
	b.mu.Unlock()
	answ, _, _ := b.sf.Do("", func() (interface{}, error) {
		defer func() { b.once = &sync.Once{} }()
		return b.v.Load(), nil
	})

	b.once.Do(func() {
		b.wait = false
		b.cond.Broadcast()
	})
	return answ
}

func NewBatcher(v *slow.Value) *Batcher {
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)
	return &Batcher{&singleflight.Group{}, v, mu, cond, &sync.Once{}, false}
}
