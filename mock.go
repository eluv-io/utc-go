package utc

// MockNowFn allows to replace the Now func variable with a mock function and returns a function to restore the default
// Now() implementation.
//
// Usage:
// defer MockNow(func() UTC { ... })()
// or
// reset := MockNow(func() UTC { ... })
// defer reset()
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
