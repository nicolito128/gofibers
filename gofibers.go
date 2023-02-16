package gofibers

import (
	"errors"
	"sync"
)

var errFiberClosed = errors.New("fiber already closed")
var errFiberNotStarted = errors.New("fiber not started")

// A Fiber Handler function
type Handler func(any, chan any)

// Fiber represents a form of corroutine handled with an interruptible function.
type Fiber struct {
	// Channel to wait the next suspend message
	lastSuspend chan any

	// If the Fiber is already started
	started bool

	// If the Fiber is already closed
	closed bool

	// Handle func
	handler Handler

	// Internal wait group
	wg *sync.WaitGroup
}

// Close the current Fiber.
func (f *Fiber) Close() error {
	if !f.closed {
		return errFiberClosed
	}

	f.closed = true
	return nil
}

// Start the Fiber execution. Receive a message that will be passed to the handler function.
func (f *Fiber) Start(msg any) error {
	if f.closed {
		return errFiberClosed
	}
	f.started = true

	go func() {
		defer close(f.lastSuspend)
		f.lastSuspend <- msg

		f.handler(msg, f.lastSuspend)
		Suspend(nil, f.wg, f.lastSuspend)
	}()

	<-f.lastSuspend
	return nil
}

// Resumes the execution of the handler from the last suspend.
func (f *Fiber) Resume() (any, error) {
	if !f.started {
		return nil, errFiberNotStarted
	}

	if f.closed {
		return nil, errFiberClosed
	}

	defer f.wg.Done()
	return <-f.lastSuspend, nil
}

// Create a new Fiber.
func New(w *sync.WaitGroup, f Handler) *Fiber {
	return &Fiber{handler: f, wg: w, lastSuspend: make(chan any)}
}

// Suspend declares an interruption in the execution of the handler.
func Suspend(val any, wg *sync.WaitGroup, response chan any) {
	wg.Add(1)
	response <- val
	wg.Wait()
}
