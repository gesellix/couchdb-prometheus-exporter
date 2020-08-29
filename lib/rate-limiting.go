package lib

import (
	"fmt"
)

type Semaphore struct {
	abort       chan struct{}
	sem         chan struct{}
	concurrency uint
}

// NewSemaphore for concurrency, concurrency == 0 means unlimited
func NewSemaphore(concurrency uint) Semaphore {
	s := Semaphore{
		sem:         make(chan struct{}, concurrency),
		abort:       make(chan struct{}), // signal to abort for goroutines waiting on a semaphore
		concurrency: concurrency,
	}
	if concurrency == 0 {
		close(s.sem) // concurrency effectively unlimited
		return s
	}
	// fill semaphore
	for i := 0; i < int(s.concurrency); i++ {
		s.sem <- struct{}{}
	}
	return s
}

// Acquire the semaphore; blocks until ready, or returns error to indicate the goroutine should abort
func (s Semaphore) Acquire() error {
	select {
	case <-s.abort:
		return fmt.Errorf("could not acquire semaphore")
	case <-s.sem:
		return nil
	}
}

func (s Semaphore) Release() {
	if s.concurrency > 0 {
		select {
		case s.sem <- struct{}{}:
		default: // should not happen unless someone double released
		}
	}
}

// Signal abort for anyone waiting on the Semaphore
func (s Semaphore) Abort() {
	select {
	case <-s.abort:
		// a check if we've already closed the abort channel
		return
	default:
		// abort channel now never blocks for receivers
		close(s.abort)
	}
}
