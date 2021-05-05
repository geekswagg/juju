// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/core/network (interfaces: ConfigSource,ConfigSourceNIC,ConfigSourceAddr)

// Package common_test is a generated GoMock package.
package common_test

import (
	net "net"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	set "github.com/juju/collections/set"
	network "github.com/juju/juju/core/network"
)

// MockConfigSource is a mock of ConfigSource interface
type MockConfigSource struct {
	ctrl     *gomock.Controller
	recorder *MockConfigSourceMockRecorder
}

// MockConfigSourceMockRecorder is the mock recorder for MockConfigSource
type MockConfigSourceMockRecorder struct {
	mock *MockConfigSource
}

// NewMockConfigSource creates a new mock instance
func NewMockConfigSource(ctrl *gomock.Controller) *MockConfigSource {
	mock := &MockConfigSource{ctrl: ctrl}
	mock.recorder = &MockConfigSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigSource) EXPECT() *MockConfigSourceMockRecorder {
	return m.recorder
}

// DefaultRoute mocks base method
func (m *MockConfigSource) DefaultRoute() (net.IP, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultRoute")
	ret0, _ := ret[0].(net.IP)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DefaultRoute indicates an expected call of DefaultRoute
func (mr *MockConfigSourceMockRecorder) DefaultRoute() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultRoute", reflect.TypeOf((*MockConfigSource)(nil).DefaultRoute))
}

// GetBridgePorts mocks base method
func (m *MockConfigSource) GetBridgePorts(arg0 string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBridgePorts", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetBridgePorts indicates an expected call of GetBridgePorts
func (mr *MockConfigSourceMockRecorder) GetBridgePorts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBridgePorts", reflect.TypeOf((*MockConfigSource)(nil).GetBridgePorts), arg0)
}

// Interfaces mocks base method
func (m *MockConfigSource) Interfaces() ([]network.ConfigSourceNIC, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Interfaces")
	ret0, _ := ret[0].([]network.ConfigSourceNIC)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Interfaces indicates an expected call of Interfaces
func (mr *MockConfigSourceMockRecorder) Interfaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Interfaces", reflect.TypeOf((*MockConfigSource)(nil).Interfaces))
}

// OvsManagedBridges mocks base method
func (m *MockConfigSource) OvsManagedBridges() (set.Strings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OvsManagedBridges")
	ret0, _ := ret[0].(set.Strings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OvsManagedBridges indicates an expected call of OvsManagedBridges
func (mr *MockConfigSourceMockRecorder) OvsManagedBridges() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OvsManagedBridges", reflect.TypeOf((*MockConfigSource)(nil).OvsManagedBridges))
}

// MockConfigSourceNIC is a mock of ConfigSourceNIC interface
type MockConfigSourceNIC struct {
	ctrl     *gomock.Controller
	recorder *MockConfigSourceNICMockRecorder
}

// MockConfigSourceNICMockRecorder is the mock recorder for MockConfigSourceNIC
type MockConfigSourceNICMockRecorder struct {
	mock *MockConfigSourceNIC
}

// NewMockConfigSourceNIC creates a new mock instance
func NewMockConfigSourceNIC(ctrl *gomock.Controller) *MockConfigSourceNIC {
	mock := &MockConfigSourceNIC{ctrl: ctrl}
	mock.recorder = &MockConfigSourceNICMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigSourceNIC) EXPECT() *MockConfigSourceNICMockRecorder {
	return m.recorder
}

// Addresses mocks base method
func (m *MockConfigSourceNIC) Addresses() ([]network.ConfigSourceAddr, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Addresses")
	ret0, _ := ret[0].([]network.ConfigSourceAddr)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Addresses indicates an expected call of Addresses
func (mr *MockConfigSourceNICMockRecorder) Addresses() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Addresses", reflect.TypeOf((*MockConfigSourceNIC)(nil).Addresses))
}

// HardwareAddr mocks base method
func (m *MockConfigSourceNIC) HardwareAddr() net.HardwareAddr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HardwareAddr")
	ret0, _ := ret[0].(net.HardwareAddr)
	return ret0
}

// HardwareAddr indicates an expected call of HardwareAddr
func (mr *MockConfigSourceNICMockRecorder) HardwareAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HardwareAddr", reflect.TypeOf((*MockConfigSourceNIC)(nil).HardwareAddr))
}

// Index mocks base method
func (m *MockConfigSourceNIC) Index() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Index")
	ret0, _ := ret[0].(int)
	return ret0
}

// Index indicates an expected call of Index
func (mr *MockConfigSourceNICMockRecorder) Index() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Index", reflect.TypeOf((*MockConfigSourceNIC)(nil).Index))
}

// IsUp mocks base method
func (m *MockConfigSourceNIC) IsUp() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUp")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsUp indicates an expected call of IsUp
func (mr *MockConfigSourceNICMockRecorder) IsUp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUp", reflect.TypeOf((*MockConfigSourceNIC)(nil).IsUp))
}

// MTU mocks base method
func (m *MockConfigSourceNIC) MTU() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MTU")
	ret0, _ := ret[0].(int)
	return ret0
}

// MTU indicates an expected call of MTU
func (mr *MockConfigSourceNICMockRecorder) MTU() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MTU", reflect.TypeOf((*MockConfigSourceNIC)(nil).MTU))
}

// Name mocks base method
func (m *MockConfigSourceNIC) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockConfigSourceNICMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockConfigSourceNIC)(nil).Name))
}

// Type mocks base method
func (m *MockConfigSourceNIC) Type() network.LinkLayerDeviceType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(network.LinkLayerDeviceType)
	return ret0
}

// Type indicates an expected call of Type
func (mr *MockConfigSourceNICMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockConfigSourceNIC)(nil).Type))
}

// MockConfigSourceAddr is a mock of ConfigSourceAddr interface
type MockConfigSourceAddr struct {
	ctrl     *gomock.Controller
	recorder *MockConfigSourceAddrMockRecorder
}

// MockConfigSourceAddrMockRecorder is the mock recorder for MockConfigSourceAddr
type MockConfigSourceAddrMockRecorder struct {
	mock *MockConfigSourceAddr
}

// NewMockConfigSourceAddr creates a new mock instance
func NewMockConfigSourceAddr(ctrl *gomock.Controller) *MockConfigSourceAddr {
	mock := &MockConfigSourceAddr{ctrl: ctrl}
	mock.recorder = &MockConfigSourceAddrMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigSourceAddr) EXPECT() *MockConfigSourceAddrMockRecorder {
	return m.recorder
}

// IP mocks base method
func (m *MockConfigSourceAddr) IP() net.IP {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IP")
	ret0, _ := ret[0].(net.IP)
	return ret0
}

// IP indicates an expected call of IP
func (mr *MockConfigSourceAddrMockRecorder) IP() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IP", reflect.TypeOf((*MockConfigSourceAddr)(nil).IP))
}

// IPNet mocks base method
func (m *MockConfigSourceAddr) IPNet() *net.IPNet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IPNet")
	ret0, _ := ret[0].(*net.IPNet)
	return ret0
}

// IPNet indicates an expected call of IPNet
func (mr *MockConfigSourceAddrMockRecorder) IPNet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IPNet", reflect.TypeOf((*MockConfigSourceAddr)(nil).IPNet))
}

// String mocks base method
func (m *MockConfigSourceAddr) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockConfigSourceAddrMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockConfigSourceAddr)(nil).String))
}
