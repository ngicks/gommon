package common

import "time"

// GetNower is getter interface of now time.Time.
//
// Use this as unexported field and mock it with ./__mock/get_now.go or any other implementation.
type GetNower interface {
	GetNow() time.Time
}

type GetNowImpl struct {
}

// GetNow implements GetNower.
func (g GetNowImpl) GetNow() time.Time {
	return time.Now()
}
