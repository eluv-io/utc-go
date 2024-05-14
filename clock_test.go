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
