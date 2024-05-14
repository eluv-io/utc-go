//go:build !race

package utc

func init() {
	setNowFn(now)
}

// clockNowFn is the function used to get the current time and can be mocked with MockNow/MockNowFn
var clock Clock

func getClock() Clock {
	ret := clock
	if ret == nil {
		ret = clockFn(now)
	}
	return ret
}

func setClock(c Clock) {
	old := clock
	clock = c

	type um interface {
		UnMocked()
	}
	if unm, ok := old.(um); ok {
		unm.UnMocked()
	}

}
