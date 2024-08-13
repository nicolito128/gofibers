package gofibers

import (
	"errors"
	"sync"
)

var (
	errFiberClosed     = errors.New("error fiber already closed")
	errFiberStarted    = errors.New("errir fiber already started")
	errFiberNotStarted = errors.New("error fiber should be started")
	errNilHandler      = errors.New("error fiber handler is nil")
)

// A fiber handler function.
//
// It receives a messages  the main goroutine.
type FiberHandler func(receive any)

// Fiber represents a form of corroutine handled with an interruptible function.
type Fiber struct {
	// Channel to wait the next suspend message
	lastSuspend chan any
	// Close signal
	closeCh chan struct{}
	// If the Fiber is already running
	started bool
	// If the Fiber is already closed
	closed bool
	// Handle func
	handler FiberHandler
	// Internal wait groups
	mainWg   *sync.WaitGroup
	threadWg *sync.WaitGroup
}

// Create a new Fiber.
func New() *Fiber {
	return &Fiber{
		mainWg:      new(sync.WaitGroup),
		lastSuspend: make(chan any, 1),
		closeCh:     make(chan struct{}),
	}
}

func (f *Fiber) Handle(handler FiberHandler) error {
	if f.closed {
		return errFiberClosed
	}

	if f.started {
		return errFiberStarted
	}

	f.handler = handler
	return nil
}

// Close the current Fiber.
func (f *Fiber) Close() error {
	if f.closed {
		return errFiberClosed
	}

	if !f.started {
		return errFiberNotStarted
	}

	f.closed = true
	close(f.lastSuspend)
	f.closeCh <- struct{}{}
	return nil
}

func (f *Fiber) Closed() <-chan struct{} {
	return f.closeCh
}

// Init the Fiber execution. Receive a message that will be passed to the handler function.
func (f *Fiber) Init(message any) error {
	if f.closed {
		return errFiberClosed
	}
	if f.handler == nil {
		return errNilHandler
	}
	f.started = true

	go func() {
		f.threadWg = new(sync.WaitGroup)
		f.handler(message)
	}()

	return nil
}

// Resumes the execution of the handler from the last suspend.
func (f *Fiber) Resume() (any, error) {
	if f.closed {
		return nil, errFiberClosed
	}

	if !f.started {
		return nil, errFiberNotStarted
	}

	v := <-f.lastSuspend
	f.threadWg.Done()
	return v, nil
}

// Suspend declares an interruption in the execution of the handler.
func (f *Fiber) Suspend(retryMessage any) {
	f.threadWg.Add(1)
	f.lastSuspend <- retryMessage
	f.threadWg.Wait()
}
