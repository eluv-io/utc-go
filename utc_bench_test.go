package utc

import "testing"

// BenchmarkNow test performances of the various now functions:
// - utc.Now using atomic
// - utc.Mono returns the straight time.Now
// - utc.WallClock strips the monotonic clock
// - utc.WallClock strips the monotonic clock and rounds to millisecond
//
// $ go test -v -bench "BenchmarkNow" -benchtime 5s -run "BenchmarkNow" ./
// goos: darwin
// goarch: amd64
// pkg: github.com/eluv-io/utc-go
// cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
// BenchmarkNow
// BenchmarkNow/utc.Now
// BenchmarkNow/utc.Now-8         	44360110	       132.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNow/utc.Mono
// BenchmarkNow/utc.Mono-8        	49817594	       122.7 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNow/utc.WallClock
// BenchmarkNow/utc.WallClock-8   	39850797	       147.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkNow/utc.WallClockMs
// BenchmarkNow/utc.WallClockMs-8 	30447565	       198.0 ns/op	       0 B/op	       0 allocs/op
// PASS
// ok  	github.com/eluv-io/utc-go	24.775s
func BenchmarkNow(b *testing.B) {
	benchmarks := []struct {
		name string
		fn   func()
	}{
		{"utc.Now", func() { _ = Now() }},
		{"utc.Mono", func() { _ = Mono.Now() }},
		{"utc.WallClock", func() { _ = WallClock.Now() }},
		{"utc.WallClockMs", func() { _ = WallClockMs.Now() }},
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
