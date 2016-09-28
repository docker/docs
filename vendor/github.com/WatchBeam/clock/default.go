package clock

import (
	"time"
)

// DefaultClock is an implementation of the Clock interface that uses standard
// time methods.
type DefaultClock struct{}

// Now returns the current local time.
func (dc DefaultClock) Now() time.Time { return time.Now() }

// After waits for the duration to elapse and then sends the current time on the returned channel.
func (dc DefaultClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

// Sleep pauses the current goroutine for at least the duration d. A negative or zero duration causes Sleep to return immediately.
func (dc DefaultClock) Sleep(d time.Duration) { time.Sleep(d) }

// Tick is a convenience wrapper for NewTicker providing access to the ticking channel only. While Tick is useful for clients that have no need to shut down the Ticker, be aware that without a way to shut it down the underlying Ticker cannot be recovered by the garbage collector; it "leaks".
func (dc DefaultClock) Tick(d time.Duration) <-chan time.Time { return time.Tick(d) }

// AfterFunc waits for the duration to elapse and then calls f in its own goroutine. It returns a Timer that can be used to cancel the call using its Stop method.
func (dc DefaultClock) AfterFunc(d time.Duration, f func()) Timer {
	return &defaultTimer{*time.AfterFunc(d, f)}
}

// NewTimer creates a new Timer that will send the current time on its channel after at least duration d.
func (dc DefaultClock) NewTimer(d time.Duration) Timer {
	return &defaultTimer{*time.NewTimer(d)}
}

// NewTicker returns a new Ticker containing a channel that will send the time with a period specified by the duration argument.
func (dc DefaultClock) NewTicker(d time.Duration) Ticker {
	return &defaultTicker{*time.NewTicker(d)}
}

type defaultTimer struct{ time.Timer }

var _ Timer = new(defaultTimer)

func (d *defaultTimer) Chan() <-chan time.Time {
	return d.C
}

type defaultTicker struct{ time.Ticker }

var _ Ticker = new(defaultTicker)

func (d *defaultTicker) Chan() <-chan time.Time {
	return d.C
}

// Default clock that uses time.Now as its time source. This is what you should
// normally use in your code.
var C = DefaultClock{}
