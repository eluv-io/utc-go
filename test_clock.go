package utc

import (
	"sync/atomic"
	"time"
)

// TestClock is a Clock that can be set to a given UTC value or reset.
// Function Now:
// - returns the previously set UTC or
// - returns the wall clock if no value or Zero was set
// A TestClock becomes effectively used as 'the global clock' after calling its
// function MockNow(). When the clock is effective, func IsMock returns true.
type TestClock struct {
	mono            bool
	millisPrecision bool
	now             *atomic.Pointer[UTC]
	isMock          *atomic.Bool
}

// NewMonoClock returns a TestClock with the monotonic clock reading.
func NewMonoClock(u ...UTC) TestClock {
	return newTestClock(true, false, u...)
}

// NewWallClock returns a TestClock with the monotonic clock reading stripped.
func NewWallClock(u ...UTC) TestClock {
	return newTestClock(false, false, u...)
}

// NewWallClockMs returns a TestClock with the monotonic clock reading stripped
// and time rounded to the millisecond.
func NewWallClockMs(u ...UTC) TestClock {
	return newTestClock(false, true, u...)
}

func newTestClock(mono, ms bool, u ...UTC) TestClock {
	ret := TestClock{
		mono:            mono,
		millisPrecision: ms,
		now:             new(atomic.Pointer[UTC]),
		isMock:          &atomic.Bool{},
	}
	if len(u) > 0 {
		ret.Set(u[0])
	}
	return ret
}

// MockNow sets this clock as the global clock. After this call utc.Now returns
// what the function Now of this clock returns.
func (c TestClock) MockNow() TestClock {
	setClock(c)
	c.isMock.Store(true)
	return c
}

func (c TestClock) unMocked() {
	c.isMock.Store(false)
}

// IsMock returns true if this clock is effectively the 'global clock'.
func (c TestClock) IsMock() bool {
	return c.isMock.Load()
}

// UnmockNow removes this clock from being the 'global clock' and resets the
// utc.Now func to the default.
func (c TestClock) UnmockNow() {
	ResetNow()
}

func (c TestClock) wc() UTC {
	if !c.mono {
		if c.millisPrecision {
			return WallClockMs()
		}
		return WallClock()
	}
	return Mono()
}

// Now returns the current time. The returned time is taken from the wall clock
// if this TestClock was started without any time or with Zero or if it was set
// to Zero or no time.
func (c TestClock) Now() UTC {
	n := c.now.Load()
	if n == nil || *n == Zero {
		return c.wc()
	}
	return *n
}

// Get returns the previously set time or Zero if it was not set.
func (c TestClock) Get() UTC {
	ret := c.now.Load()
	if ret == nil {
		return Zero
	}
	return *ret
}

// Set sets the given UTC time and returns the previously set time or Zero if
// it was not set. If the parameter is Zero, the function works as if Unset was
// called.
func (c TestClock) Set(u UTC) UTC {
	return c.set(u)
}

// Unset unsets the current time and returns the previously set value.
// This resets this TestClock to behave as a wall clock in future calls to Now.
// This is equivalent to calling c.Set(Zero)
func (c TestClock) Unset() UTC {
	return c.set(Zero)
}

func (c TestClock) set(u UTC) UTC {
	var n *UTC

	switch u {
	case Zero:
		n = nil
	default:
		w := u
		if !c.mono {
			w = w.StripMono()
			if c.millisPrecision {
				w = w.Round(time.Millisecond)
			}
		}
		n = &w
	}

	ret := c.now.Swap(n)
	if ret == nil {
		return Zero
	}
	return *ret
}

// Add adds the given duration to the UTC time of this TestClock and returns the
// resulting UTC.
// If this TestClock was started without a time, the addition is made on top of
// the current wall clock (and results in a time in the future when t is positive).
func (c TestClock) Add(t time.Duration) UTC {
	n := c.Now()
	ret := n.Add(t)
	c.Set(ret)
	return ret
}

// SetNow sets this TestClock to the current wall clock and returns the previously
// set value. This is equivalent to calling c.Set( wall-clock ).
func (c TestClock) SetNow() UTC {
	return c.Set(c.wc())
}
