package utc

import (
	"sync/atomic"
	"time"
)

// Chrono is a Clock that can be set to a given UTC value or reset.
// Function Now:
// - returns the previously set UTC or
// - returns the wall clock if no value or Zero was set
type Chrono struct {
	mono bool
	ms   bool
	u    *atomic.Pointer[UTC]
}

// NewWatch returns a Chrono where the monotonic clock reading has been stripped
func NewWatch(u ...UTC) Chrono {
	return newChrono(false, false, u...)
}

// NewWatchMs returns a Chrono where the monotonic clock reading has been stripped
// and times are rounded to the millisecond. This is useful in tests where times
// have to be serialised and compared.
func NewWatchMs(u ...UTC) Chrono {
	return newChrono(false, true, u...)
}

// NewChrono returns a Chrono that uses the monotonic clock.
func NewChrono(u ...UTC) Chrono {
	return newChrono(true, false, u...)
}

func newChrono(mono, ms bool, u ...UTC) Chrono {
	ret := Chrono{
		mono: mono,
		ms:   ms,
		u:    &atomic.Pointer[UTC]{},
	}
	ret.Set(u...)
	return ret
}

func (c Chrono) wc() UTC {
	if !c.mono {
		if c.ms {
			return WallClockMs.Now()
		}
		return WallClock.Now()
	}
	return Mono.Now()
}

// Now returns the current time. The returned time is taken from the wall clock
// if this Chrono was started without any time or with Zero.
func (c Chrono) Now() UTC {
	n := c.u.Load()
	if n == nil || *n == Zero {
		return c.wc()
	}
	return *n
}

// Get returns the previously set time or Zero if it was not set.
func (c Chrono) Get() UTC {
	ret := c.u.Load()
	if ret == nil {
		return Zero
	}
	return *ret
}

// Set sets the given UTC time and returns the previously set time or Zero if
// it was not set. If no parameter is specified, the current time is unset and
// calls to this Chrono.Now will be returning the actual current time from the
// wall clock.
func (c Chrono) Set(u ...UTC) UTC {
	var n *UTC
	if len(u) > 0 {
		switch u[0] {
		case Zero:
			n = nil
		default:
			w := u[0]
			if !c.mono {
				w = w.StripMono()
				if c.ms {
					w = w.Round(time.Millisecond)
				}
			}
			n = &w
		}
	}
	ret := c.u.Swap(n)
	if ret == nil {
		return Zero
	}
	return *ret
}

// Add adds the given duration to the UTC time of this Chrono and returns the
// resulting UTC.
// If this Chrono was started without a time, the addition is made on top of the
// current wall clock (and results in a time in the future when t is positive).
func (c Chrono) Add(t time.Duration) UTC {
	n := c.Now()
	ret := n.Add(t)
	c.Set(ret)
	return ret
}

// Wall reset this Chrono to behave as a wall clock in future calls to Now and
// returns the previously set value. This is equivalent to calling c.Set(Zero)
func (c Chrono) Wall() UTC {
	return c.Set(Zero)
}

// SetNow sets this Chrono to the current wall clock and returns the previously
// set value. This is equivalent to calling c.Set( wall-clock ).
func (c Chrono) SetNow() UTC {
	return c.Set(c.wc())
}
