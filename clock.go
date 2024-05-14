package utc

import "time"

type Clock interface {
	Now() UTC
}

var (
	// WallClock is a Clock returning the wall clock (i.e. where the monotonic
	// clock reading has been stripped).
	WallClock = wallClock{}

	// WallClockMs is a WallClock rounded to the millisecond.
	WallClockMs = wallClock{ms: true}

	// Mono is a Clock returning the current time with the monotonic clock.
	Mono = mono{}
)

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
