// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/worker/caasunitprovisioner (interfaces: ProvisioningStatusSetter)

// Package caasunitprovisioner is a generated GoMock package.
package caasunitprovisioner

import (
	gomock "github.com/golang/mock/gomock"
	status "github.com/juju/juju/core/status"
	reflect "reflect"
)

// MockProvisioningStatusSetter is a mock of ProvisioningStatusSetter interface
type MockProvisioningStatusSetter struct {
	ctrl     *gomock.Controller
	recorder *MockProvisioningStatusSetterMockRecorder
}

// MockProvisioningStatusSetterMockRecorder is the mock recorder for MockProvisioningStatusSetter
type MockProvisioningStatusSetterMockRecorder struct {
	mock *MockProvisioningStatusSetter
}

// NewMockProvisioningStatusSetter creates a new mock instance
func NewMockProvisioningStatusSetter(ctrl *gomock.Controller) *MockProvisioningStatusSetter {
	mock := &MockProvisioningStatusSetter{ctrl: ctrl}
	mock.recorder = &MockProvisioningStatusSetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProvisioningStatusSetter) EXPECT() *MockProvisioningStatusSetterMockRecorder {
	return m.recorder
}

// SetOperatorStatus mocks base method
func (m *MockProvisioningStatusSetter) SetOperatorStatus(arg0 string, arg1 status.Status, arg2 string, arg3 map[string]interface{}) error {
	ret := m.ctrl.Call(m, "SetOperatorStatus", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetOperatorStatus indicates an expected call of SetOperatorStatus
func (mr *MockProvisioningStatusSetterMockRecorder) SetOperatorStatus(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOperatorStatus", reflect.TypeOf((*MockProvisioningStatusSetter)(nil).SetOperatorStatus), arg0, arg1, arg2, arg3)
}
