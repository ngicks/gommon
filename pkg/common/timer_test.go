package common_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ngicks/gommon/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimerRealStop(t *testing.T) {
	t.Parallel()

	timer := common.NewTimerReal()

	assertNotExpired := func() {
		t.Helper()
		select {
		case <-timer.C():
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
	case <-timer.C():
		t.Fatal()
	}

	// Making sure that Reset(0) and time.Sleep combo really expires the timer.
	timer.Reset(0)
	time.Sleep(time.Microsecond)
	select {
	case <-timer.C():
	default:
		t.Fatal()
	}

	timer.Reset(0)
	time.Sleep(time.Microsecond)
	assert.False(t, timer.Stop())
	<-timer.C()
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
			<-timer.C()
			then := time.Now()
			require.GreaterOrEqual(t, int64(then.Sub(now)), int64(tt))
		})
	}
}
