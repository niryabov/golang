//go:build !solution

package cond

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

// Cond implements a condition variable, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
//
// Each Cond has an associated Locker L (often a *sync.Mutex or *sync.RWMutex),
// which must be held when changing the condition and
// when calling the Wait method.
type Cond struct {
	L     Locker
	cntCh chan int
	lCh   chan int
}

// New returns a new Cond with Locker l.
func New(l Locker) *Cond {

	c := &Cond{L: l, cntCh: make(chan int, 1), lCh: make(chan int, 1)}
	c.cntCh <- 0

	return c
}

// Wait atomically unlocks c.L and suspends execution
// of the calling goroutine. After later resuming execution,
// Wait locks c.L before returning. Unlike in other systems,
// Wait cannot return unless awoken by Broadcast or Signal.

// Because c.L is not locked when Wait first resumes, the caller
// typically cannot assume that the condition is true when
// Wait returns. Instead, the caller should Wait in a loop:
//
//	c.L.Lock()
//	for !condition() {
//	    c.Wait()
//	}
//	... make use of condition ...
//	c.L.Unlock()
func (c *Cond) Wait() {
	cnt := <-c.cntCh
	c.cntCh <- (cnt + 1)
	c.L.Unlock()

	<-c.lCh
	c.L.Lock()
}

// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Signal() {
	cnt := <-c.cntCh
	if cnt == 0 {
		c.cntCh <- 0
		return
	}
	c.cntCh <- cnt - 1
	c.lCh <- 1

}

// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Broadcast() {

	cnt := <-c.cntCh
	c.cntCh <- cnt
	for i := 0; i < cnt; i++ {
		go c.Signal()
	}

}
