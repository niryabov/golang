//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	tokens   chan bool
	stopped  bool
	exit     chan bool
	interval time.Duration
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	tokens := make(chan bool, maxCount)

	for i := 0; i < maxCount; i++ {
		tokens <- true
	}
	return &Limiter{tokens, false, make(chan bool, 1), interval}

}

func delay(l *Limiter, tick <-chan time.Time) {
	for {
		select {
		case <-l.exit:
			l.exit <- true
			return
		case <-tick:
			l.tokens <- true
			return
		default:
		}
	}
}

func (l *Limiter) Acquire(ctx context.Context) error {

	for {
		if l.stopped {
			l.exit <- true
			defer func() { <-l.exit }()
			return ErrStopped
		}
		select {

		case <-ctx.Done():
			l.exit <- true
			defer func() { <-l.exit }()

			return ctx.Err()
		case <-l.tokens:
			if l.interval.Nanoseconds() != 0 {
				tick := time.Tick(l.interval)
				go delay(l, tick)
			} else {
				l.tokens <- true
			}
			return nil
		default:
		}

	}
}

func (l *Limiter) Stop() {
	l.exit <- true
	l.stopped = true
	time.Sleep(time.Millisecond)

	defer func() { <-l.exit }()
}
