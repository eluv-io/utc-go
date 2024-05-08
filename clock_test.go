package utc

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWallClock(t *testing.T) {
	now := WallClock.Now()
	ms := WallClockMs.Now()
	require.True(t, ms.Sub(now) <= time.Millisecond)
}

func TestWallClockMs(t *testing.T) {
	now := WallClockMs.Now()
	bb, err := json.Marshal(now)
	require.NoError(t, err)

	n := UTC{}
	err = json.Unmarshal(bb, &n)
	require.NoError(t, err)
	require.Equal(t, now, n)
}

func TestMono(t *testing.T) {
	now := Mono.Now()
	bb, err := json.Marshal(now)
	require.NoError(t, err)

	n := UTC{}
	err = json.Unmarshal(bb, &n)
	require.NoError(t, err)
	require.NotEqual(t, now, n)
}

func TestWatch(t *testing.T) {
	now := now()
	wh := NewWatch(now)
	ch := NewChrono(now)

	for _, c := range []Chrono{wh, ch} {
		oneSecAfter := c.Add(time.Second)
		require.Equal(t, time.Second, oneSecAfter.Sub(now))
		require.Equal(t, time.Second, c.Now().Sub(now))
	}
}

func TestChrono(t *testing.T) {
	now := now()
	w := NewWatch(now)
	c := NewChrono(now)

	require.Equal(t, w.Now().UnixMilli(), c.Now().UnixMilli())
	require.Equal(t, time.Duration(0), w.Now().Sub(c.Now()))

	// reset to follow the wall clock
	n := c.Wall()
	require.Equal(t, now, n)
	time.Sleep(time.Millisecond * 200)
	n2 := c.Now()
	d := n2.Sub(n)
	require.True(t, d >= time.Millisecond*200 && d < time.Millisecond*250)
	time.Sleep(time.Millisecond * 200)
	n2 = c.Now()
	d = n2.Sub(n)
	require.True(t, d >= time.Millisecond*400 && d < time.Millisecond*450)

	// set to the current time
	n2 = c.SetNow()
	require.Equal(t, Zero, n2)
	n2 = c.Now()
	time.Sleep(time.Millisecond * 200)
	n3 := c.Now()
	require.Equal(t, n2, n3)
}

func TestChronoGet(t *testing.T) {
	now := now()
	m := NewChrono(now)
	require.Equal(t, now, m.Get())

	w := NewWatch(now)
	require.Equal(t, now.StripMono(), w.Get())

	ws := NewWatchMs(now)
	require.Equal(t, now.StripMono().Round(time.Millisecond), ws.Get())
}
