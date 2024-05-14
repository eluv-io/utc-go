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
	mono bool
	ms   bool
	u    *atomic.Pointer[UTC]
	m    *atomic.Bool
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
		mono: mono,
		ms:   ms,
		u:    new(atomic.Pointer[UTC]),
		m:    &atomic.Bool{},
	}
	ret.Set(u...)
	return ret
}

// MockNow sets this clock as the global clock. After this call utc.Now returns
// what the function Now of this clock returns.
func (c TestClock) MockNow() TestClock {
	setClock(c)
	c.m.Store(true)
	return c
}

func (c TestClock) unMocked() {
	c.m.Store(false)
}

// IsMock returns true if this clock is effectively the 'global clock'.
func (c TestClock) IsMock() bool {
	return c.m.Load()
}

// Reset removes this clock from being the 'global clock' and resets the utc.Now
// func to the default.
func (c TestClock) Reset() {
	ResetNow()
}

func (c TestClock) wc() UTC {
	if !c.mono {
		if c.ms {
			return WallClockMs.Now()
		}
		return WallClock.Now()
	}
	return Mono.Now()
}

// Now returns the current time. The returned time is taken from the wall clock
// if this TestClock was started without any time or with Zero or if it was set
// to Zero or no time.
func (c TestClock) Now() UTC {
	n := c.u.Load()
	if n == nil || *n == Zero {
		return c.wc()
	}
	return *n
}

// Get returns the previously set time or Zero if it was not set.
func (c TestClock) Get() UTC {
	ret := c.u.Load()
	if ret == nil {
		return Zero
	}
	return *ret
}

// Set sets the given UTC time and returns the previously set time or Zero if
// it was not set. If no parameter is specified, the current time is unset and
// calls to this TestClock.Now will be returning the actual current time from
// the wall clock.
func (c TestClock) Set(u ...UTC) UTC {
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

// Wall reset this TestClock to behave as a wall clock in future calls to Now
// and returns the previously set value. This is equivalent to calling c.Set(Zero)
func (c TestClock) Wall() UTC {
	return c.Set(Zero)
}

// SetNow sets this TestClock to the current wall clock and returns the previously
// set value. This is equivalent to calling c.Set( wall-clock ).
func (c TestClock) SetNow() UTC {
	return c.Set(c.wc())
}
