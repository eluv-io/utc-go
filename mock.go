package utc

// MockNowFn allows to replace the Now func variable with a mock function and returns a function to restore the default
// Now() implementation.
//
// Usage:
// defer MockNow(func() UTC { ... })()
// or
// reset := MockNow(func() UTC { ... })
// defer reset()
//
// Using MockNowFn like below and updating the local variable from multiple
// go routines is racy, hence the function is marked as deprecated.
//
//	now := utc.Now()
//	reset := utc.MockNowFn(func() utc.UTC { return now })
//	...
//	now = now.Add(time.Duration(i+1)*time.Second))
//
// Deprecated use one of the NewXXClock().MockNow()
func MockNowFn(fn func() UTC) (restore func()) {
	setNowFn(fn)
	return ResetNow
}

// MockNow allows to replace the Now func variable with a function that returns
// the given constant time and returns itself a function to restore the default
// Now() implementation.
//
// Usage:
// defer MockNow(utc.MustParse("2020-01-01"))()
// or
// reset := MockNow(utc.MustParse("2020-01-01"))
// defer reset()
func MockNow(time UTC) (restore func()) {
	return MockNowFn(func() UTC {
		return time
	})
}

// MockNowClock mocks now with a test clock. Equivalent to calling clock.MockNow().
func MockNowClock(clock TestClock) {
	clock.MockNow()
}
