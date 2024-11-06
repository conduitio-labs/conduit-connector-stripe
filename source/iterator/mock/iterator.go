// Code generated by MockGen. DO NOT EDIT.
// Source: iterator.go
//
// Generated by this command:
//
//	mockgen -package mock -source iterator.go -destination ./mock/iterator.go
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	models "github.com/conduitio-labs/conduit-connector-stripe/models"
	gomock "go.uber.org/mock/gomock"
)

// MockStripe is a mock of Stripe interface.
type MockStripe struct {
	ctrl     *gomock.Controller
	recorder *MockStripeMockRecorder
}

// MockStripeMockRecorder is the mock recorder for MockStripe.
type MockStripeMockRecorder struct {
	mock *MockStripe
}

// NewMockStripe creates a new mock instance.
func NewMockStripe(ctrl *gomock.Controller) *MockStripe {
	mock := &MockStripe{ctrl: ctrl}
	mock.recorder = &MockStripeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStripe) EXPECT() *MockStripeMockRecorder {
	return m.recorder
}

// GetEvent mocks base method.
func (m *MockStripe) GetEvent(createdAt int64, startingAfter, endingBefore string) (models.EventResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEvent", createdAt, startingAfter, endingBefore)
	ret0, _ := ret[0].(models.EventResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEvent indicates an expected call of GetEvent.
func (mr *MockStripeMockRecorder) GetEvent(createdAt, startingAfter, endingBefore any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEvent", reflect.TypeOf((*MockStripe)(nil).GetEvent), createdAt, startingAfter, endingBefore)
}

// GetResource mocks base method.
func (m *MockStripe) GetResource(arg0 string) (models.ResourceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResource", arg0)
	ret0, _ := ret[0].(models.ResourceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResource indicates an expected call of GetResource.
func (mr *MockStripeMockRecorder) GetResource(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResource", reflect.TypeOf((*MockStripe)(nil).GetResource), arg0)
}
