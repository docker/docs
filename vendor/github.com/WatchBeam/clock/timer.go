package clock

import (
	"sync"
	"time"
)

type mockTimer struct {
	c       chan time.Time
	release chan bool

	mu     sync.Mutex
	clock  Clock
	active bool
	target time.Time
}

var _ Timer = new(mockTimer)

func (m *mockTimer) getTarget() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.target
}

func (m *mockTimer) setInactive() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// If a release was sent in the meantime, that means a new timer
	// was started or that we already stopped manually
	select {
	case <-m.release:
		return
	default:
		m.active = false
	}

}

func (m *mockTimer) wait() {
	select {
	case <-m.clock.After(m.target.Sub(m.clock.Now())):
		m.c <- m.clock.Now()
		m.setInactive()
	case <-m.release:
	}
}

func (m *mockTimer) Chan() <-chan time.Time {
	return m.c
}

func (m *mockTimer) Reset(d time.Duration) bool {
	var wasActive bool
	m.mu.Lock()
	defer m.mu.Unlock()

	wasActive, m.active = m.active, true
	m.target = m.clock.Now().Add(d)

	if wasActive {
		m.release <- true
	}
	go m.wait()

	return wasActive
}

func (m *mockTimer) Stop() bool {
	var wasActive bool
	m.mu.Lock()
	defer m.mu.Unlock()

	wasActive, m.active = m.active, false
	if wasActive {
		m.release <- true
	}

	return wasActive
}

// Creates a new Timer using the provided Clock. You should not use this
// directly outside of unit tests; use Clock.NewTimer().
func NewMockTimer(c Clock) Timer {
	return &mockTimer{
		c:       make(chan time.Time, 1),
		release: make(chan bool),
		clock:   c,
	}
}
