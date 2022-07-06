package gommon_test

import (
	"testing"
	"time"

	"github.com/ngicks/gommon"
	"github.com/stretchr/testify/require"
)

func TestTimer(t *testing.T) {
	timer := gommon.NewTimerImpl()

	now := time.Now()
	timer.Reset(now, now)

	require.InDelta(t, now.UnixNano(), (<-timer.Channel()).UnixNano(), float64(time.Millisecond))

	timer.Reset(now.Add(-time.Second), now)
	// emit fast.
	require.InDelta(t, now.UnixNano(), (<-timer.Channel()).UnixNano(), float64(time.Millisecond))
}
