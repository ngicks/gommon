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
	type testCase struct {
		d time.Duration
	}

	cases := []testCase{
		{
			d: time.Millisecond,
		},
		{
			d: time.Microsecond,
		},
		{
			d: time.Second,
		},
		{
			d: time.Second + time.Millisecond*50,
		},
		{
			d: 2 * time.Second,
		},
	}

	for idx, testCase := range cases {
		tt := testCase
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			t.Parallel()

			timer := common.NewTimerReal()

			now := time.Now()
			timer.Reset(tt.d)
			<-timer.Channel()
			then := time.Now()
			require.GreaterOrEqual(t, int64(then.Sub(now)), int64(tt.d))
		})
	}
}

func TestTimerRealResetTo(t *testing.T) {
	type testCase struct {
		to time.Time
	}

	now := time.Now()
	cases := []testCase{
		{
			to: now.Add(time.Millisecond),
		},
		{
			to: now.Add(time.Microsecond),
		},
		{
			to: now.Add(time.Second),
		},
		{
			to: now.Add(time.Second + time.Millisecond*50),
		},
		{
			to: now.Add(2 * time.Second),
		},
	}

	for idx, testCase := range cases {
		tt := testCase
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			t.Parallel()

			timer := common.NewTimerReal()

			timer.ResetTo(tt.to)
			<-timer.Channel()
			then := time.Now()
			require.InDelta(t, then.UnixNano(), tt.to.UnixNano(), float64(10*time.Millisecond))
		})
	}
}
