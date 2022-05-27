// Code generated by MockGen. DO NOT EDIT.
// Source: source/iterator/interface.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	sdk "github.com/conduitio/conduit-connector-sdk"
	models "github.com/conduitio/conduit-connector-stripe/models"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Next mocks base method.
func (m *MockRepository) Next() (sdk.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next")
	ret0, _ := ret[0].(sdk.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Next indicates an expected call of Next.
func (mr *MockRepositoryMockRecorder) Next() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockRepository)(nil).Next))
}

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
func (mr *MockStripeMockRecorder) GetEvent(createdAt, startingAfter, endingBefore interface{}) *gomock.Call {
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
func (mr *MockStripeMockRecorder) GetResource(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResource", reflect.TypeOf((*MockStripe)(nil).GetResource), arg0)
}
