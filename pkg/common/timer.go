package common

//go:generate mockgen -source timer.go -destination __mock/timer.go

import "time"

// ITimer is timer interface.
//
// Use this as unexported field and mock it with ./__mock/timer.go or any other implementation.
type ITimer interface {
	// Channel is equivalent of timer.C
	Channel() <-chan time.Time
	// Reset changes the timer to expire after duration d.
	Reset(d time.Duration)
	// Reset changes the timer to expire at time to.
	ResetTo(to time.Time)
	// Stop stops this timer.
	Stop()
}

var _ ITimer = NewTimerImpl()

// TimerImpl is a struct that implements ITimer.
type TimerImpl struct {
	GetNower
	*time.Timer
}

// NewTimerImpl returns newly created TimerImpl.
// Unlike time.NewTimer, this creates stopped timer.
func NewTimerImpl() *TimerImpl {
	timer := time.NewTimer(time.Second)
	if !timer.Stop() {
		<-timer.C
	}
	return &TimerImpl{
		GetNower: GetNowImpl{},
		Timer:    timer,
	}
}

func (t *TimerImpl) Channel() <-chan time.Time {
	return t.C
}

func (t *TimerImpl) Stop() {
	if !t.Timer.Stop() {
		// non-blocking receive.
		// in case of racy concurrent receivers.
		select {
		case <-t.C:
		default:
		}
	}
}

func (t *TimerImpl) Reset(d time.Duration) {
	t.Stop()
	t.Timer.Reset(d)
}

func (t *TimerImpl) ResetTo(to time.Time) {
	t.Reset(to.Sub(t.GetNow()))
}
