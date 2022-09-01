// Code generated by MockGen. DO NOT EDIT.
// Source: get_now.go

// Package mock_common is a generated GoMock package.
package mock_common

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockGetNower is a mock of GetNower interface.
type MockGetNower struct {
	ctrl     *gomock.Controller
	recorder *MockGetNowerMockRecorder
}

// MockGetNowerMockRecorder is the mock recorder for MockGetNower.
type MockGetNowerMockRecorder struct {
	mock *MockGetNower
}

// NewMockGetNower creates a new mock instance.
func NewMockGetNower(ctrl *gomock.Controller) *MockGetNower {
	mock := &MockGetNower{ctrl: ctrl}
	mock.recorder = &MockGetNowerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetNower) EXPECT() *MockGetNowerMockRecorder {
	return m.recorder
}

// GetNow mocks base method.
func (m *MockGetNower) GetNow() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNow")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetNow indicates an expected call of GetNow.
func (mr *MockGetNowerMockRecorder) GetNow() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNow", reflect.TypeOf((*MockGetNower)(nil).GetNow))
}