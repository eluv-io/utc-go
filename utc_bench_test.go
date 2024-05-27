package utc

import (
	"testing"
	"time"
)

// BenchmarkNow test performances of the various now functions:
// - utc.Now using nowFn
// - utc.Mono returns the straight time.Now
// - utc.WallClock strips the monotonic clock
// - utc.WallClockMs strips the monotonic clock and rounds to millisecond
//
// gaillard:utc-go $ go test -v -bench "BenchmarkNow" -benchtime 5s -run "BenchmarkNow" ./
// goos: darwin
// goarch: amd64
// pkg: github.com/eluv-io/utc-go
// cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
// BenchmarkNow
// BenchmarkNow/utc.Now
// BenchmarkNow/utc.Now-8         	50482977	       118.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNow/utc.Mono
// BenchmarkNow/utc.Mono-8        	50519804	       118.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNow/utc.WallClock
// BenchmarkNow/utc.WallClock-8   	39826977	       145.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNow/utc.WallClockMs
// BenchmarkNow/utc.WallClockMs-8 	30610243	       195.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNowClock
// BenchmarkNowClock/utc.Now
// BenchmarkNowClock/utc.Now-8    	42252702	       141.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNowClock/utc.Mono
// BenchmarkNowClock/utc.Mono-8   	50646658	       118.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNowClock/utc.WallClock
// BenchmarkNowClock/utc.WallClock-8         	41166252	       146.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNowClock/utc.WallClockMs
// BenchmarkNowClock/utc.WallClockMs-8       	30708422	       195.6 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	github.com/eluv-io/utc-go	49.199s
func BenchmarkNow(b *testing.B) {
	doBenchmarkNow(b, false)
}

func BenchmarkNowClock(b *testing.B) {
	doBenchmarkNow(b, true)
}

func doBenchmarkNow(b *testing.B, withClock bool) {
	b.StopTimer()
	if !withClock {
		nowFn = now
	} else {
		allowClock()
	}
	b.StartTimer()
	benchmarks := []struct {
		name string
		fn   func()
	}{
		{"utc.Now", func() { _ = Now() }},
		{"utc.Mono", func() { _ = Mono() }},
		{"utc.WallClock", func() { _ = WallClock() }},
		{"utc.WallClockMs", func() { _ = WallClockMs() }},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				bm.fn()
			}
		})
	}
}

// BenchmarkTimeNow
// BenchmarkTimeNow/time.Now
// BenchmarkTimeNow/time.Now-8         	 9869420	       111.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkTimeNow/utc.Now
// BenchmarkTimeNow/utc.Now-8          	 9474093	       125.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkTimeNow/utc.monotonic
// BenchmarkTimeNow/utc.monotonic-8    	 9015828	       128.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkTimeNow/utc.Mono
// BenchmarkTimeNow/utc.Mono-8         	 9543199	       127.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkTimeNow/utc.WallClock
// BenchmarkTimeNow/utc.WallClock-8    	 7693312	       153.0 ns/op	       0 B/op	       0 allocs/op
// BenchmarkTimeNow/utc.WallClockMs
// BenchmarkTimeNow/utc.WallClockMs-8  	 5173959	       223.3 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTimeNow(b *testing.B) {
	b.StopTimer()
	nowFn = now
	b.StartTimer()
	benchmarks := []struct {
		name string
		fn   func()
	}{
		{"time.Now", func() { _ = time.Now() }},
		{"utc.Now", func() { _ = Now() }},
		{"utc.monotonic", func() { _ = monotonic.Now() }},
		{"utc.Mono", func() { _ = Mono() }},
		{"utc.WallClock", func() { _ = WallClock() }},
		{"utc.WallClockMs", func() { _ = WallClockMs() }},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				bm.fn()
			}
		})
	}
}
