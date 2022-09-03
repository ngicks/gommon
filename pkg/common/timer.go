package common

//go:generate mockgen -source timer.go -destination __mock/timer.go

import "time"

// Timer is an interface equivalent to time.Timer. This is an interface so that it can be mocked.
//
// Use this as an unexported field and swap out in tests.const
// In non-test env, TimerReal should suffice. in tests, mock is pre-generated in ./__mock/timer.go, by mockgen.
type Timer interface {
	// Channel is equivalent of timer.C
	Channel() <-chan time.Time
	// Reset changes the timer to expire after duration d. If timer is already expired Channel will be drained before reset.
	Reset(d time.Duration)
	// Reset changes the timer to expire at time to. If timer is already expired Channel will be drained.
	ResetTo(to time.Time)
	// Stop stops this timer. If timer is already expired Channel will be drained.
	Stop()
}

var _ Timer = NewTimerReal()

// TimerReal is a struct that implements Timer.
//
// TimerReal uses real(runtime) timer.
type TimerReal struct {
	*time.Timer
}

// NewTimerReal returns newly created TimerReal.
// This creates stopped timer unlike time.NewTimer.
func NewTimerReal() *TimerReal {
	timer := time.NewTimer(time.Second)
	if !timer.Stop() {
		<-timer.C
	}
	return &TimerReal{
		Timer: timer,
	}
}

func (t *TimerReal) Channel() <-chan time.Time {
	return t.C
}

func (t *TimerReal) Stop() {
	if !t.Timer.Stop() {
		// non-blocking receive.
		// in case of racy concurrent receivers.
		select {
		case <-t.C:
		default:
		}
	}
}

func (t *TimerReal) Reset(d time.Duration) {
	t.Stop()
	t.Timer.Reset(d)
}

func (t *TimerReal) ResetTo(to time.Time) {
	t.Reset(time.Until(to))
}
