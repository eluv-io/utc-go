package utc

import "time"

type Clock interface {
	Now() UTC
}

// ClockFn is a function implementing Clock
type ClockFn func() UTC

func (fn ClockFn) Now() UTC {
	return fn()
}

var (
	wall      = wallClock{}
	wallMs    = wallClock{ms: true}
	monotonic = mono{}
)

// WallClock returns the wall clock (i.e. where the monotonic clock reading has been stripped).
func WallClock() UTC {
	return wall.Now()
}

// WallClockMs is like WallClock rounded to the millisecond.
func WallClockMs() UTC {
	return wallMs.Now()
}

// Mono returns the current time with the monotonic clock.
func Mono() UTC {
	return monotonic.Now()
}

type wallClock struct {
	ms bool
}

func (w wallClock) Now() UTC {
	// strip mono
	ret := New(now().Time.Truncate(0))
	if w.ms {
		ret = ret.Round(time.Millisecond)
	}
	return ret
}

type mono struct{}

func (mono) Now() UTC {
	return now()
}
