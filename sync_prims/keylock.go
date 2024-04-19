//go:build !solution

package keylock

import (
	"sync"
)

type KeyLock struct {
	mu *sync.Mutex
	m  map[string]chan bool
	// locking chan bool
}

func New() *KeyLock {
	// ch := make(chan bool, 1)
	// ch <- true
	return &KeyLock{&sync.Mutex{}, make(map[string]chan bool)}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
external:
	for {
		taken_keys := make([]string, 0)

		unlock = func() {
			for _, key := range taken_keys {
				l.m[key] <- true
			}

		}

		l.mu.Lock()
		for _, key := range keys {
			_, ok := l.m[key]
			if !ok {
				l.m[key] = make(chan bool, 1)
				l.m[key] <- true
			}

			select {
			case <-cancel:
				unlock()
				l.mu.Unlock()
				return true, unlock
			case <-l.m[key]:
				taken_keys = append(taken_keys, key)
			default:
				unlock()
				l.mu.Unlock()
				continue external
			}

		}
		l.mu.Unlock()
		return false, unlock

	}
}
