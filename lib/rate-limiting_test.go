package lib

import (
	"testing"
	"time"
)

// worker goroutine; should take no less than 1ms to complete.
func worker(sem Semaphore, results chan struct{}) {
	err := sem.Acquire()
	if err != nil {
		return
	}
	time.Sleep(1 * time.Millisecond)
	results <- struct{}{}
	sem.Release()
}

func TestSemaphore(t *testing.T) {
	// Test1: test concurrent workers, limited to 10 at a time.
	sem := NewSemaphore(10)
	// channel buffered to receive all results
	results := make(chan struct{}, 100)

	for i := 0; i < 100; i++ {
		go worker(sem, results)
	}
	time.Sleep(20 * time.Millisecond) // plenty of time
	if len(results) != 100 {
		t.Error("Workers did not complete in time")
	}
}

func TestUnlimitedSemaphore(t *testing.T) {
	// Test2: test unlimited concurrent workers
	sem := NewSemaphore(0)
	results := make(chan struct{}, 100)

	for i := 0; i < 100; i++ {
		go worker(sem, results)
	}
	time.Sleep(20 * time.Millisecond) // plenty of time
	if len(results) != 100 {
		t.Error("Workers did not complete in time")
	}
}

func TestSerialSemaphore(t *testing.T) {
	// Test3: test serialised workers, ie. 1 at a time;
	sem := NewSemaphore(1)
	results := make(chan struct{}, 100)
	err := sem.Acquire()
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 100; i++ {
		go worker(sem, results)
	}
	sem.Release()
	time.Sleep(20 * time.Millisecond)
	if len(results) > 20 { // should never complete this many in the allotted time
		t.Error("Workers completed job too fast", len(results))
	}
}

func TestAbort(t *testing.T) {
	// Test abort functionality
	sem := NewSemaphore(1)
	results := make(chan struct{}, 100)
	for i := 0; i < 100; i++ {
		go worker(sem, results)
	}
	// let workers make some progress
	time.Sleep(10 * time.Millisecond)
	sem.Abort()
	time.Sleep(200 * time.Millisecond)
	// there should never be 100 results
	if len(results) == 100 {
		t.Error("Somehow all workers completed their jobs despite an abort")
	}
}
