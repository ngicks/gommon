package common

import "time"

// NowGetter is getter interface of now time.Time.
//
// Use this as unexported field and mock it with ./__mock/get_now.go or any other implementation.
type NowGetter interface {
	GetNow() time.Time
}

type NowGetterReal struct {
}

// GetNow implements NowGetter.
func (g NowGetterReal) GetNow() time.Time {
	return time.Now()
}
