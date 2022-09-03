package common_test

import (
	"testing"
	"time"

	"github.com/ngicks/gommon/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestTimer(t *testing.T) {
	timer := common.NewTimerReal()

	now := time.Now()
	timer.ResetTo(now)

	require.InDelta(t, now.UnixNano(), (<-timer.Channel()).UnixNano(), float64(time.Millisecond))

	timer.ResetTo(now.Add(-time.Second))
	// emit fast.
	require.InDelta(t, now.UnixNano(), (<-timer.Channel()).UnixNano(), float64(time.Millisecond))
}
