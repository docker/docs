package tracer

import (
	"fmt"
	"github.com/stretchr/stew/strings"
	"time"
)

// deferredTrace holds the format string and args sent for deferred processing
type deferredTrace struct {
	fmt   string
	args  []interface{}
	level Level
	now   time.Time
}

// Trace is a type that contains an actual trace
type Trace struct {
	// data stores the trace string
	Data string
	// level stores the level of this trace
	Level Level
	// timestamp stores the time at which this trace occurred
	TimeStamp time.Time
}

// tracer is a type providing trace capabilities to a go application
type Tracer struct {
	// data stores the individual string entries in the tracer
	data []Trace
	// level determines what level a trace must be to actually be logged
	level Level
	// deferred holds whether processing is deferred or now
	deferred bool
	// deferredData holds the deferredTrace for the traces to be processed
	deferredData []deferredTrace
}

type Level int

const (
	// LevelEverything defines a trace level that is guarenteed to trace everything.  This trace level
	// should never be passed to the Trace command.
	LevelEverything Level = iota
	// LevelDebug defines a trace level that should be used for verbose tracing.
	LevelDebug
	// LevelInfo defines a trace level that should be used for normal activity tracing.
	LevelInfo
	// LevelWarning defines a trace level that should be used for tracing warnings.
	LevelWarning
	// LevelError defines a trace level that should be used for error tracing.
	LevelError
	// LevelCritical defines a trace level that should be used for critical error tracing.
	LevelCritical
	// LevelNothing defines a trace level that is guarenteed to trace nothing.  This trace level
	// should never be passed to the Trace command.
	LevelNothing
)

// New creates a new tracer
//
// The level argument is used to filter the trace to the desired level of detail.
// For example, a trace level of LevelEverything will log everything, where a trace level of LevelWarning
// will log only warnings, errors and criticals.
func New(level Level) *Tracer {

	tracer := new(Tracer)

	tracer.level = level
	tracer.data = make([]Trace, 0, 100)

	return tracer

}

// NewDeferred creates a new tracer that defers processing until Process() is called
//
// The level argument is used to filter the trace to the desired level of detail.
// For example, a trace level of LevelEverything will log everything, where a trace level of LevelWarning
// will log only warnings, errors and criticals.
func NewDeferred(level Level) *Tracer {

	tracer := new(Tracer)

	tracer.level = level
	tracer.data = make([]Trace, 0, 100)
	tracer.deferred = true
	tracer.deferredData = make([]deferredTrace, 0, 100)

	return tracer

}

// Level gets the current level of this Tracer.
func (t *Tracer) Level() Level {

	if t == nil {
		return LevelNothing
	}

	return t.level
}

// Should gets whether the tracer should trace at the specified level.
func (t *Tracer) Should(level Level) bool {
	return t.Level() <= level
}

// Trace saves a piece of trace data at the current time.
func (t *Tracer) Trace(level Level, format string, args ...interface{}) {

	if t == nil {
		return
	}

	if level >= t.level && level < LevelNothing {
		if t.deferred {
			t.deferredData = append(t.deferredData, deferredTrace{fmt: format, args: args, level: level, now: time.Now()})
		} else {
			trace := Trace{fmt.Sprintf(format, args...), level, time.Now()}
			t.data = append(t.data, trace)
		}
	} else if level <= LevelEverything {
		panic("tracer: level is invalid: Cannot Trace with LevelEverything or below.")
	} else if level >= LevelNothing {
		panic("tracer: level is invalid: Cannot Trace with LevelNothing or above.")
	}

}

// TraceDebug is a helper method for tracing at the LevelDebug level
func (t *Tracer) TraceDebug(format string, args ...interface{}) {
	t.Trace(LevelDebug, format, args...)
}

// TraceInfo is a helper method for tracing at the LevelInfo level
func (t *Tracer) TraceInfo(format string, args ...interface{}) {
	t.Trace(LevelInfo, format, args...)
}

// TraceWarning is a helper method for tracing at the LevelWarning level
func (t *Tracer) TraceWarning(format string, args ...interface{}) {
	t.Trace(LevelWarning, format, args...)
}

// TraceError is a helper method for tracing at the LevelError level
func (t *Tracer) TraceError(format string, args ...interface{}) {
	t.Trace(LevelError, format, args...)
}

// TraceCritical is a helper method for tracing at the LevelCritical level
func (t *Tracer) TraceCritical(format string, args ...interface{}) {
	t.Trace(LevelCritical, format, args...)
}

// Returns a copy of the trace data
func (t *Tracer) Data() []Trace {

	if t == nil {
		return nil
	}

	copiedTraces := make([]Trace, len(t.data))

	copy(copiedTraces, t.data)

	return copiedTraces

}

//Process processes deferred traces
func (t *Tracer) Process() {
	if t.deferred && len(t.deferredData) > 0 {
		for _, defTrace := range t.deferredData {
			trace := Trace{fmt.Sprintf(defTrace.fmt, defTrace.args...), defTrace.level, defTrace.now}
			t.data = append(t.data, trace)
		}
		t.deferredData = make([]deferredTrace, 0, 100)
	}
}

// Returns a copy of the trace data as an array of string
func (t *Tracer) StringData() []string {

	if t == nil {
		return nil
	}

	t.Process()

	stringTraces := make([]string, 0, len(t.data))

	for _, trace := range t.data {
		stringTraces = append(stringTraces, fmt.Sprintf("TRACE: %s\t%s\t\t%s", trace.TimeStamp.String(), LevelToString(trace.Level), trace.Data))
	}

	return stringTraces

}

// String gets a nicely formatted string of the trace data.
func (t *Tracer) String() string {

	if t == nil {
		return ""
	}

	t.Process()

	return strings.MergeStrings("\n", strings.JoinStrings("\n", t.StringData()...))
}

// Returns a copy of the trace data, filtered by trace level
func (t *Tracer) Filter(level Level) []Trace {

	if t == nil {
		return nil
	}

	t.Process()

	filteredTraces := make([]Trace, 0, len(t.data))

	for _, trace := range t.data {

		if trace.Level == level {
			filteredTraces = append(filteredTraces, trace)
		}

	}

	return filteredTraces

}

// Returns a string representation of the level
func LevelToString(level Level) string {

	switch level {
	case LevelEverything:
		return "LevelEverything"
	case LevelDebug:
		return "     LevelDebug"
	case LevelInfo:
		return "      LevelInfo"
	case LevelWarning:
		return "   LevelWarning"
	case LevelError:
		return "     LevelError"
	case LevelCritical:
		return "  LevelCritical"
	case LevelNothing:
		return "   LevelNothing"
	}
	return ""
}
