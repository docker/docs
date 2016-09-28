package clock

import (
	"time"
)

// The Clock interface provides time-based functionality. It should be used
// rather than the `time` package in situations where you want to mock things.
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) Timer
	NewTimer(d time.Duration) Timer
	NewTicker(d time.Duration) Ticker
}

// The Timer is an interface for time.Timer, and can also be swapped in mocks.
// This *does* change its API so that it can fit into an interface -- rather
// than using the channel at .C, you should call Chan() and use the
// returned channel just as you would .C.
type Timer interface {
	Chan() <-chan time.Time
	Reset(d time.Duration) bool
	Stop() bool
}

// The Timer is an interface for time.Ticket, and can also be swapped in mocks.
// This *does* change its API so that it can fit into an interface -- rather
// than using the channel at .C, you should call Chan() and use the
// returned channel just as you would .C.
type Ticker interface {
	Chan() <-chan time.Time
	Stop()
}
