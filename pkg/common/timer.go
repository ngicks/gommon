package common

import (
	"sync"
	"time"
)

// Timer is an interface equivalent to time.Timer. This is an interface so that it can be mocked.
// Unlike time.Timer, newly created Timer should be stopped timer. Return value of Reset is omitted as
// this interface does not need to care about compatibility with old programs.
//
// Use this as an unexported field and swap out in tests.
// In non-test env, TimerReal should suffice. in tests, use FakeTimer or mock which is pre-generated in ./__mock/timer.go by mockgen.
type Timer interface {
	// C is equivalent of timer.C
	C() <-chan time.Time
	// Reset changes the timer to expire after duration d. If the timer is already expired, C is drained before reset.
	Reset(d time.Duration)
	// Stop prevents timer from firing. It returns true if it successfully stop the timer, false if it has already expired or been stopped.
	Stop() bool
}

var _ Timer = NewTimerReal()

// TimerReal implements Timer using a real(runtime) time.Timer.
type TimerReal struct {
	T *time.Timer
}

// NewTimerReal returns newly created TimerReal.
// This creates stopped timer unlike time.NewTimer.
func NewTimerReal() *TimerReal {
	timer := time.NewTimer(time.Second)
	if !timer.Stop() {
		<-timer.C
	}
	return &TimerReal{
		T: timer,
	}
}

func (t *TimerReal) C() <-chan time.Time {
	return t.T.C
}

func (t *TimerReal) Stop() bool {
	return t.T.Stop()
}

func (t *TimerReal) Reset(d time.Duration) {
	if !t.Stop() {
		// non-blocking receive.
		// in case of racy concurrent receivers.
		select {
		case <-t.T.C:
		default:
		}
	}
	t.T.Reset(d)
}

var _ Timer = NewTimerFake()

// TimerFake implements Timer to swap TimerReal in test files.
type TimerFake struct {
	sync.Mutex
	Channel   chan time.Time
	ResetArgs []*time.Duration   // nil means Stop call.
	ResetCh   chan time.Duration // can synchronize with Reset if call on channel-receive operator on this channel before a Rest call.
	Expired   bool
}

func NewTimerFake() *TimerFake {
	return &TimerFake{
		ResetArgs: make([]*time.Duration, 0),
		Channel:   make(chan time.Time),
		ResetCh:   make(chan time.Duration, 1),
	}
}

func (t *TimerFake) C() <-chan time.Time {
	return t.Channel
}

func (t *TimerFake) Reset(d time.Duration) {
	t.Lock()
	t.Expired = false
	t.ResetArgs = append(t.ResetArgs, &d)
	t.Unlock()

	select {
	case t.ResetCh <- d:
	default:
	}
}

func (t *TimerFake) Stop() bool {
	t.Lock()
	defer t.Unlock()
	t.ResetArgs = append(t.ResetArgs, nil)
	return !t.Expired
}

func (t *TimerFake) Send(tt time.Time) {
	t.Lock()
	t.Expired = true
	t.Unlock()

	t.Channel <- tt
}

func (t *TimerFake) SetExpired(expired bool) {
	t.Lock()
	defer t.Unlock()

	t.Expired = expired
}

func (t *TimerFake) ExhaustResetCh() {
	for {
		select {
		case <-t.ResetCh:
		default:
			return
		}
	}
}
