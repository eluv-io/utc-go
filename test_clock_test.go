package utc_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eluv-io/utc-go"
)

func TestNewWallClock(t *testing.T) {
	u := utc.Now()
	mono := utc.NewMonoClock(u)
	wall := utc.NewWallClock(u)
	wms := utc.NewWallClockMs(u)

	require.Equal(t, mono.Get().StripMono(), wall.Get())
	require.Equal(t, wall.Get().Round(time.Millisecond), wms.Get())
}

func TestClockMock(t *testing.T) {
	wc := utc.NewWallClock().MockNow()
	assert.True(t, wc.IsMock())
	defer wc.UnmockNow()
	assert.False(t, utc.Now().Equal(utc.Zero))

	d2020 := utc.MustParse("2020-01-01")
	wc.Set(d2020)
	assert.True(t, utc.Now().Equal(d2020))
	curr := wc.Get()
	assert.Equal(t, d2020, curr)

	prev := wc.Add(time.Hour)
	assert.Equal(t, d2020.Add(time.Hour), prev)
	assert.Equal(t, utc.Now(), wc.Now())
	assert.Equal(t, d2020.Add(time.Hour), wc.Now())

	prev = wc.Set(utc.Zero)
	assert.Equal(t, d2020.Add(time.Hour), prev)

	assertRestored := func() {
		timeNow := time.Now()
		utcNow := utc.Now()
		assert.True(t, timeNow.Before(utcNow.UTC()) || timeNow.Equal(utcNow.UTC()), "time.Now=%s utc.Now=%s", timeNow, utcNow)
	}
	assertRestored()

	wc.Set(d2020.Add(time.Hour))
	assert.Equal(t, utc.Now(), wc.Now())
	wc.UnmockNow()
	assert.False(t, wc.IsMock())
	assertRestored()
	assert.NotEqual(t, utc.Now(), wc.Now())
}

func TestClockSince(t *testing.T) {
	nowUTC := utc.Now()
	wc := utc.NewWallClock(nowUTC).MockNow()
	wc.Add(time.Minute)
	defer wc.UnmockNow()

	assert.Equal(t, time.Minute, utc.Since(nowUTC))
}

func TestClockUntil(t *testing.T) {
	thenUTC := utc.Now().Add(time.Second)
	wc := utc.NewWallClock(thenUTC).MockNow()
	wc.Add(-time.Minute)
	defer wc.UnmockNow()

	assert.Equal(t, time.Minute, utc.Until(thenUTC))
}

// TestRace shows a race between utc.Now() and assigning nowFn = nowFnClock
// and how to avoid it.
func TestRace(t *testing.T) {
	// without the call below to ResetNow, running the test with -race fails. Note
	// that installing any clock - like with utc.MockNowClock(utc.NewMonoClock()) -
	// instead of calling ResetNow also works.
	utc.ResetNow()

	stopCh := make(chan struct{})
	defer close(stopCh)

	go func() {
		for {
			select {
			case <-stopCh:
				break
			default:
				utc.Now()
				time.Sleep(time.Millisecond)
			}
		}
	}()
	for i := 0; i < 100; i++ {
		utc.NewWallClock().MockNow()
	}

}
