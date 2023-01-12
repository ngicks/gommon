package timing

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/ngicks/gommon/pkg/common"
)

// swap out this if tests need to.
var timerFactory = func() common.Timer {
	return common.NewTimerReal()
}

type pollParam struct {
	ctx context.Context
}

func newPollParam() *pollParam {
	return &pollParam{
		ctx: context.Background(),
	}
}

type pollOption func(*pollParam)

func SetPollContext(ctx context.Context) pollOption {
	return func(pp *pollParam) {
		pp.ctx = ctx
	}
}

// PollUntil polls on the predicate function until the predicate returns true.
// The predicate is called at intervals of interval.
// If the predicate does not return true after timeout, PollUntil will return false ok.
func PollUntil(predicate func(ctx context.Context) bool, interval time.Duration, timeout time.Duration, options ...pollOption) (ok bool) {
	param := newPollParam()

	for _, opt := range options {
		opt(param)
	}

	ctx := param.ctx

	predCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	doneCh := make(chan struct{})

	defer func() {
		<-doneCh
	}()

	var called atomic.Bool
	done := func() {
		if called.CompareAndSwap(false, true) {
			cancel()
			close(doneCh)
		}
	}

	wait := make(chan struct{})
	defer func() {
		<-wait
	}()

	go func() {
		defer func() { close(wait) }()
		t := timerFactory()
		defer t.Stop()
		for {
			select {
			case <-doneCh:
				return
			default:
			}
			if predicate(predCtx) {
				break
			}
			t.Reset(interval)
			select {
			case <-t.C():
			case <-doneCh:
				return
			}
		}
		done()
	}()

	t := timerFactory()
	t.Reset(timeout)
	defer t.Stop()
	select {
	case <-ctx.Done():
		done()
		return false
	case <-t.C():
		done()
		return false
	default:
		select {
		case <-ctx.Done():
			done()
			return false
		case <-t.C():
			done()
			return false
		case <-doneCh:
			return true
		}
	}
}
