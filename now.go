package utc

import (
	"time"
)

// nowFn is the function used to get the current time and can be mocked with
// MockNow, MockNowFn, MockNowClock as described in comments of setClock and allowClock.
var nowFn func() UTC

func init() {
	nowFn = now
}

// New creates a new UTC instance from the given time. Use utc.Now() to get the
// current time.
func New(t time.Time) UTC {
	return UTC{Time: t.UTC(), mono: t}
}

// Now returns the current time as UTC instance. Now can be mocked for tests: see MockNow() function.
func Now() UTC {
	return nowFn()
}

// WallNow returns Now as a wall clock, i.e. with the monotonic clock reading stripped.
// WallNow is equivalent to calling Now().StripMono().
func WallNow() UTC {
	return Now().StripMono()
}

// WallNowMs returns Now as a wall clock rounded to the millisecond.
// WallNowMs is equivalent to calling WallNow().Round(time.Millisecond) and useful in tests where UTC instances are
// serialized and compared.
func WallNowMs() UTC {
	return WallNow().Round(time.Millisecond)
}

// now is the default, non-mocked value of Now.
func now() UTC {
	return New(time.Now())
}

// ResetNow resets the Now func to the default implementation.
func ResetNow() {
	setNowFn(now)
}

// setNowFn sets the given function as the Clock to use for tests.
func setNowFn(fn func() UTC) {
	setClock(ClockFn(fn))
}
