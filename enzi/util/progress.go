package util

import (
	"github.com/docker/distribution/context"
)

// ProgressLoggingIterator allows for iterating over a slice. The slice
// can be looped over like so:
//
//     for progress.Next() {
//         index := progress.Index()
//         ... do something with index ...
//     }
//
// It should print a log message at certain times: 0%, 100%, every 1000
// iterations, etc.
type ProgressLoggingIterator interface {
	Next() bool
	Index() int
}

// ProgressLogger is a simple wrapper around a logger which makes it easy to
// log percentage progress at 0%, 100%, and every N iterations.
type ProgressLogger struct {
	Logger   context.Logger
	Message  string
	StepSize int
	Length   int

	started    bool
	current    int
	nextThresh int
}

var _ ProgressLoggingIterator = (*ProgressLogger)(nil)

// start starts this progress logger. A "0% Complete" log message is emitted.
func (pl *ProgressLogger) start() {
	pl.started = true
	pl.current = 0
	pl.nextThresh = pl.StepSize
	pl.logPercent(0.0)
}

// complete completes this progress logger. A "100% Complete" log message is
// emitted.
func (pl *ProgressLogger) complete() {
	if !pl.started {
		// Not yet started or already complete.
		return
	}

	pl.started = false // Reset the progress.
	pl.logPercent(100.0)
}

// Next increments the current value and may produce a log message if the
// current value crosses a progress threshold. The progress logger is started
// if it has not already been. Returns false when complete.
func (pl *ProgressLogger) Next() bool {
	// Check if this is the start.
	if !pl.started {
		pl.start()

		if pl.Length == 0 {
			// Nothing to progress over.
			pl.complete()
			return false
		}

		return true
	}

	pl.current++

	// Check for 100% completion. Handle this explicitly due to rounding
	// errors.
	if pl.current == pl.Length {
		pl.complete()
		return false
	}

	if pl.current == pl.nextThresh {
		// We've crossed a step threshhold.
		pl.nextThresh += pl.StepSize
		currentProgress := float64(pl.current) / float64(pl.Length)
		pl.logPercent(currentProgress * 100.0)
	}

	return true
}

// Index returns the current position of this progress logger.
func (pl *ProgressLogger) Index() int {
	return pl.current
}

func (pl *ProgressLogger) logPercent(percent float64) {
	pl.Logger.Infof("%s - %.2f%% Complete", pl.Message, percent)
}
