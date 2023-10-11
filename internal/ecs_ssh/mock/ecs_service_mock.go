// Code generated by MockGen. DO NOT EDIT.
// Source: ecs_service.go
//
// Generated by this command:
//
//	mockgen -source=ecs_service.go -package=mock -destination=./mock/ecs_service_mock.go
//
// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockECSService is a mock of ECSService interface.
type MockECSService struct {
	ctrl     *gomock.Controller
	recorder *MockECSServiceMockRecorder
}

// MockECSServiceMockRecorder is the mock recorder for MockECSService.
type MockECSServiceMockRecorder struct {
	mock *MockECSService
}

// NewMockECSService creates a new mock instance.
func NewMockECSService(ctrl *gomock.Controller) *MockECSService {
	mock := &MockECSService{ctrl: ctrl}
	mock.recorder = &MockECSServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockECSService) EXPECT() *MockECSServiceMockRecorder {
	return m.recorder
}

// GetTargetIDByServiceName mocks base method.
func (m *MockECSService) GetTargetIDByServiceName(ctx context.Context, clusterName, serviceName, containerName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTargetIDByServiceName", ctx, clusterName, serviceName, containerName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTargetIDByServiceName indicates an expected call of GetTargetIDByServiceName.
func (mr *MockECSServiceMockRecorder) GetTargetIDByServiceName(ctx, clusterName, serviceName, containerName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTargetIDByServiceName", reflect.TypeOf((*MockECSService)(nil).GetTargetIDByServiceName), ctx, clusterName, serviceName, containerName)
}

// GetTargetIDByTaskID mocks base method.
func (m *MockECSService) GetTargetIDByTaskID(ctx context.Context, clusterName, taskID, containerName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTargetIDByTaskID", ctx, clusterName, taskID, containerName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTargetIDByTaskID indicates an expected call of GetTargetIDByTaskID.
func (mr *MockECSServiceMockRecorder) GetTargetIDByTaskID(ctx, clusterName, taskID, containerName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTargetIDByTaskID", reflect.TypeOf((*MockECSService)(nil).GetTargetIDByTaskID), ctx, clusterName, taskID, containerName)
}