package utc

import (
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
)

// functions in this file implement use of a Clock to return the current UTC.
// This is intended for tests and won't happen when running production code.
// See comments of allowClock.

// clocker wraps a Clock
type clocker struct {
	c Clock
}

var (
	// nowFnMu protects access to nowFn when nowFn is set to nowFnClock. This is
	// intended to be used during tests only.
	nowFnMu sync.Mutex
	// atomicClock stores the current Clock
	atomicClock atomic.Pointer[clocker]
)

// nowFnClock is the function used to get the current time via a Clock during tests.
func nowFnClock() UTC {
	return getClock().Now()
}

// allowClock allows using a custom Clock and is called during test whenever setClock is called.
//   - in production code nowFn is initialised to function now, is never updated and
//     therefore does not need any protection against races.
//   - in test code nowFn is replaced with nowFnClock which will incur a lookup into atomicClock
//
// The function minimizes changes to nowFn, hence - although the mutex does not
// protect read access to nowFn - the probability of race is reduced. An example of
// such a race will happen whenever a test that does not mock UTC has a goroutine
// that does not terminate (and keep calling utc.Now) while another test starts and
// installs a clock to mock UTC (see TestRace in test_clock_test.go). This kind of
// case clearly results from a faulty test and the -race flag helps in detecting it.
// On the other hand, a package that would like to prevent any such race could implement
// an initializer using:
// //go:build race
//
//	func init(){
//	 utc.ResetNow()
//	}
func allowClock() {
	nowFnMu.Lock()
	defer nowFnMu.Unlock()
	if runtime.FuncForPC(reflect.ValueOf(nowFn).Pointer()) == runtime.FuncForPC(reflect.ValueOf(nowFnClock).Pointer()) {
		//nowFnClock already set
		return
	}
	nowFn = nowFnClock
}

// getClock returns the current Clock stored in atomicClock or the default 'now'
// function if no clock was stored.
func getClock() Clock {
	var fn Clock = ClockFn(now)
	m := atomicClock.Load()
	if m != nil && m.c != nil {
		fn = m.c
	}
	return fn
}

// setClock sets c to be the current clock. This function is intended to be used
// in tests only through one of the MockNowXyz functions. It first calls allowClock,
// then replaces the current clock with c.
func setClock(c Clock) {
	allowClock()

	old := getClock()
	n := &clocker{c: c}
	atomicClock.Store(n)

	type um interface {
		unMocked()
	}
	if unm, ok := old.(um); ok {
		// notify the previous clock that it is no more the current 'mock'
		// unMocked is currently implemented only by TestClock
		unm.unMocked()
	}
}
