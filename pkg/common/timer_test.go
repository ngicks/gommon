package common_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ngicks/gommon/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestTimerRealStop(t *testing.T) {
	t.Parallel()

	timer := common.NewTimerReal()

	assertNotExpired := func() {
		select {
		case <-timer.Channel():
			t.Fatal()
		default:
		}
	}

	assertNotExpired()

	timer.Reset(time.Second)
	assertNotExpired()
	timer.Stop()
	assertNotExpired()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
	case <-timer.Channel():
		t.Fatal()
	}

	// Making sure that Reset(0) and time.Sleep combo really expires the timer.
	timer.Reset(0)
	time.Sleep(time.Microsecond)
	select {
	case <-timer.Channel():
	default:
		t.Fatal()
	}

	timer.Reset(0)
	time.Sleep(time.Microsecond)
	timer.Stop()
	assertNotExpired()
}

func TestTimerRealReset(t *testing.T) {
	t.Parallel()

	cases := []time.Duration{
		time.Millisecond,
		time.Microsecond,
		time.Second,
		time.Second + time.Millisecond*50,
		2 * time.Second,
	}

	for idx, testCase := range cases {
		tt := testCase
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			t.Parallel()

			timer := common.NewTimerReal()

			now := time.Now()
			timer.Reset(tt)
			<-timer.Channel()
			then := time.Now()
			require.GreaterOrEqual(t, int64(then.Sub(now)), int64(tt))
		})
	}
}

func TestTimerRealResetTo(t *testing.T) {
	t.Parallel()

	now := time.Now()
	cases := []time.Time{
		now.Add(time.Millisecond),
		now.Add(time.Microsecond),
		now.Add(time.Second),
		now.Add(time.Second + time.Millisecond*50),
		now.Add(2 * time.Second),
	}

	for idx, testCase := range cases {
		tt := testCase
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			t.Parallel()

			timer := common.NewTimerReal()

			timer.ResetTo(tt)
			<-timer.Channel()
			then := time.Now()
			require.InDelta(t, then.UnixNano(), tt.UnixNano(), float64(10*time.Millisecond))
		})
	}
}
