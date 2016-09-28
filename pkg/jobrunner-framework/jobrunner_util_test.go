package jobrunner

import (
	"fmt"
	"time"
)

type Waiter struct {
	ch   chan struct{}
	name string
}

func NewWaiter(name string) *Waiter {
	return &Waiter{
		ch:   make(chan struct{}),
		name: name,
	}
}

func (w *Waiter) Wait(timeout time.Duration) error {
	timeoutCh := time.After(timeout)
	select {
	case <-w.ch:
		return nil
	case <-timeoutCh:
		return fmt.Errorf("Timed out waiting on waiter '%s' after %s", w.name, timeout)
	}
	panic("this shouldn't happen")
}

func (w *Waiter) Close() error {
	close(w.ch)
	return nil
}

func ChWait(name string, ch <-chan struct{}, timeout time.Duration) error {
	timeoutCh := time.After(timeout)
	select {
	case <-ch:
		return nil
	case <-timeoutCh:
		return fmt.Errorf("Timed out waiting for '%s' after %s", name, timeout)
	}
	panic("this shouldn't happen")
}
