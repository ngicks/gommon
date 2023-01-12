package common

import "time"

// Timer is an interface equivalent to time.Timer. This is an interface so that it can be mocked.
// Unlike time.Timer, newly created Timer should be stopped timer. Return value of Reset is omitted as
// this interface does not need to care about compatibility with old programs.
//
// Use this as an unexported field and swap out in tests.
// In non-test env, TimerReal should suffice. in tests, mock is pre-generated in ./__mock/timer.go, by mockgen.
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

func (t *TimerReal) Reset(d time.Duration) {
	t.Stop()
	t.Timer.Reset(d)
}

func (t *TimerReal) ResetTo(to time.Time) {
	t.Reset(time.Until(to))
}
