package clock

import (
	"time"
)

type mockTicker struct {
	c    chan time.Time
	stop chan bool

	clock    Clock
	interval time.Duration
	start    time.Time
}

var _ Ticker = new(mockTicker)

// note: this probably does not function the same way as the time.Timer
// in the event that the clock skips more than the timer interval. I've
// not yet dug deep into the runtimeTimer to see how that works.
// PRs are appreciated!
func (m *mockTicker) wait() {
	for i := time.Duration(1); true; i++ {
		delta := m.start.Add(m.interval * i).Sub(m.clock.Now())

		select {
		case <-m.stop:
			return
		case <-m.clock.After(delta):
			m.c <- time.Now()
		}
	}
}

func (m *mockTicker) Chan() <-chan time.Time {
	return m.c
}

func (m *mockTicker) Stop() {
	m.stop <- true
}

// Creates a new Ticker using the provided Clock. You should not use this
// directly outside of unit tests; use Clock.NewTicker().
func NewMockTicker(c Clock, interval time.Duration) Ticker {
	t := &mockTicker{
		c:        make(chan time.Time),
		stop:     make(chan bool),
		interval: interval,
		start:    c.Now(),
		clock:    c,
	}
	go t.wait()

	return t
}
