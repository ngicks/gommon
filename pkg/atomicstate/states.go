// Code generated by github.com/ngicks/gommon/cmd/generate_state_impl/generate_state_impl.go. DO NOT EDIT.
package atomicstate

import "sync/atomic"


// WorkingState is atomic state primitive.
// It holds a boolean state corresponds to its name.
type WorkingState struct {
	s uint32
}

// IsWorking is atomic state checker.
// It returns true if state is set, and vice versa.
func (s *WorkingState) IsWorking() bool {
	return atomic.LoadUint32(&s.s) == 1
}

// SetWorking is atomic state setter.
// It tries to set its state based on to.
// If first element of to is false, it tries set it to false,
// true otherwise.
//
// swapped is true when it is successfully set, false if it is already set to the state.
func (s *WorkingState) SetWorking(to ...bool) (swapped bool) {
	setTo := true
	if len(to) > 0 {
		setTo = to[0]
	}

	if setTo {
		return atomic.CompareAndSwapUint32(&s.s, 0, 1)
	} else {
		return atomic.CompareAndSwapUint32(&s.s, 1, 0)
	}
}

// NewWorkingState builds splitted WorkingState wrapper.
// Either or both can be embedded and/or used as unexported member to hide its setter.
func NewWorkingState() (checker *WorkingStateChecker, setter *WorkingStateSetter) {
	s := new(WorkingState)
	checker = &WorkingStateChecker{s}
	setter = &WorkingStateSetter{s}
	return
}

// WorkingStateSetter is sipmle wrapper of WorkingState.
// It only exposes IsWorking.
type WorkingStateChecker struct {
	s *WorkingState
}

// IsWorking is atomic state checker.
// It returns true if state is set, and vice versa.
func (s *WorkingStateChecker) IsWorking() bool {
	return s.s.IsWorking()
}

// WorkingStateSetter is sipmle wrapper of WorkingState.
// It only exposes SetWorking. 
type WorkingStateSetter struct {
	s *WorkingState
}

// SetWorking is atomic state setter.
// It tries to set its state based on to.
// If first element of to is false, it tries set it to false,
// true otherwise.
//
// swapped is true when it is successfully set, false if it is already set to the state.
func (s *WorkingStateSetter) SetWorking(to ...bool) (swapped bool) {
	return s.s.SetWorking(to...)
}

// EndedState is atomic state primitive.
// It holds a boolean state corresponds to its name.
type EndedState struct {
	s uint32
}

// IsEnded is atomic state checker.
// It returns true if state is set, and vice versa.
func (s *EndedState) IsEnded() bool {
	return atomic.LoadUint32(&s.s) == 1
}

// SetEnded is atomic state setter.
// It tries to set its state based on to.
// If first element of to is false, it tries set it to false,
// true otherwise.
//
// swapped is true when it is successfully set, false if it is already set to the state.
func (s *EndedState) SetEnded(to ...bool) (swapped bool) {
	setTo := true
	if len(to) > 0 {
		setTo = to[0]
	}

	if setTo {
		return atomic.CompareAndSwapUint32(&s.s, 0, 1)
	} else {
		return atomic.CompareAndSwapUint32(&s.s, 1, 0)
	}
}

// NewEndedState builds splitted EndedState wrapper.
// Either or both can be embedded and/or used as unexported member to hide its setter.
func NewEndedState() (checker *EndedStateChecker, setter *EndedStateSetter) {
	s := new(EndedState)
	checker = &EndedStateChecker{s}
	setter = &EndedStateSetter{s}
	return
}

// EndedStateSetter is sipmle wrapper of EndedState.
// It only exposes IsEnded.
type EndedStateChecker struct {
	s *EndedState
}

// IsEnded is atomic state checker.
// It returns true if state is set, and vice versa.
func (s *EndedStateChecker) IsEnded() bool {
	return s.s.IsEnded()
}

// EndedStateSetter is sipmle wrapper of EndedState.
// It only exposes SetEnded. 
type EndedStateSetter struct {
	s *EndedState
}

// SetEnded is atomic state setter.
// It tries to set its state based on to.
// If first element of to is false, it tries set it to false,
// true otherwise.
//
// swapped is true when it is successfully set, false if it is already set to the state.
func (s *EndedStateSetter) SetEnded(to ...bool) (swapped bool) {
	return s.s.SetEnded(to...)
}