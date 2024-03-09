// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aws/aws-application-networking-k8s/pkg/deploy/lattice (interfaces: ListenerManager)

// Package lattice is a generated GoMock package.
package lattice

import (
	context "context"
	reflect "reflect"

	lattice "github.com/aws/aws-application-networking-k8s/pkg/model/lattice"
	vpclattice "github.com/aws/aws-sdk-go/service/vpclattice"
	gomock "github.com/golang/mock/gomock"
)

// MockListenerManager is a mock of ListenerManager interface.
type MockListenerManager struct {
	ctrl     *gomock.Controller
	recorder *MockListenerManagerMockRecorder
}

// MockListenerManagerMockRecorder is the mock recorder for MockListenerManager.
type MockListenerManagerMockRecorder struct {
	mock *MockListenerManager
}

// NewMockListenerManager creates a new mock instance.
func NewMockListenerManager(ctrl *gomock.Controller) *MockListenerManager {
	mock := &MockListenerManager{ctrl: ctrl}
	mock.recorder = &MockListenerManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockListenerManager) EXPECT() *MockListenerManagerMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockListenerManager) Delete(arg0 context.Context, arg1 *lattice.Listener) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockListenerManagerMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockListenerManager)(nil).Delete), arg0, arg1)
}

// List mocks base method.
func (m *MockListenerManager) List(arg0 context.Context, arg1 string) ([]*vpclattice.ListenerSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*vpclattice.ListenerSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockListenerManagerMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockListenerManager)(nil).List), arg0, arg1)
}

// Upsert mocks base method.
func (m *MockListenerManager) Upsert(arg0 context.Context, arg1 *lattice.Listener, arg2 *lattice.Service, arg3 *lattice.TargetGroup) (lattice.ListenerStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(lattice.ListenerStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upsert indicates an expected call of Upsert.
func (mr *MockListenerManagerMockRecorder) Upsert(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockListenerManager)(nil).Upsert), arg0, arg1, arg2, arg3)
}
