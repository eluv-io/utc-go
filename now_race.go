//go:build race

package utc

import "sync/atomic"

type clocker struct {
	c Clock
}

var atomicClock atomic.Pointer[clocker]

func getClock() Clock {
	var fn Clock = clockFn(now)
	m := atomicClock.Load()
	if m != nil && m.c != nil {
		fn = m.c
	}
	return fn
}

func setClock(c Clock) {
	allowClock()

	old := getClock()
	n := &clocker{c: c}
	atomicClock.Store(n)

	type um interface {
		unMocked()
	}
	if unm, ok := old.(um); ok {
		unm.unMocked()
	}
}
