// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go
//
// Generated by this command:
//
//	mockgen -source=manager.go -package=mock -destination=./mock/manager_mock.go
//
// Package mock is a generated GoMock package.
package mock

import (
	sessions "aws-ssh/internal/sessions"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// ExecSession mocks base method.
func (m *MockManager) ExecSession(in *sessions.ExecInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecSession", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExecSession indicates an expected call of ExecSession.
func (mr *MockManagerMockRecorder) ExecSession(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecSession", reflect.TypeOf((*MockManager)(nil).ExecSession), in)
}

// PortForwardingSession mocks base method.
func (m *MockManager) PortForwardingSession(in *sessions.PortForwardingInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PortForwardingSession", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// PortForwardingSession indicates an expected call of PortForwardingSession.
func (mr *MockManagerMockRecorder) PortForwardingSession(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PortForwardingSession", reflect.TypeOf((*MockManager)(nil).PortForwardingSession), in)
}

// PortForwardingToRemoteHostSession mocks base method.
func (m *MockManager) PortForwardingToRemoteHostSession(in *sessions.PortForwardingToRemoteInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PortForwardingToRemoteHostSession", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// PortForwardingToRemoteHostSession indicates an expected call of PortForwardingToRemoteHostSession.
func (mr *MockManagerMockRecorder) PortForwardingToRemoteHostSession(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PortForwardingToRemoteHostSession", reflect.TypeOf((*MockManager)(nil).PortForwardingToRemoteHostSession), in)
}
