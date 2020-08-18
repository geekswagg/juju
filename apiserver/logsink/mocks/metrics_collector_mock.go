// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/logsink (interfaces: MetricsCollector)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	prometheus "github.com/prometheus/client_golang/prometheus"
	reflect "reflect"
)

// MockMetricsCollector is a mock of MetricsCollector interface
type MockMetricsCollector struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsCollectorMockRecorder
}

// MockMetricsCollectorMockRecorder is the mock recorder for MockMetricsCollector
type MockMetricsCollectorMockRecorder struct {
	mock *MockMetricsCollector
}

// NewMockMetricsCollector creates a new mock instance
func NewMockMetricsCollector(ctrl *gomock.Controller) *MockMetricsCollector {
	mock := &MockMetricsCollector{ctrl: ctrl}
	mock.recorder = &MockMetricsCollectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMetricsCollector) EXPECT() *MockMetricsCollectorMockRecorder {
	return m.recorder
}

// Connections mocks base method
func (m *MockMetricsCollector) Connections() prometheus.Gauge {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connections")
	ret0, _ := ret[0].(prometheus.Gauge)
	return ret0
}

// Connections indicates an expected call of Connections
func (mr *MockMetricsCollectorMockRecorder) Connections() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connections", reflect.TypeOf((*MockMetricsCollector)(nil).Connections))
}

// LogReadCount mocks base method
func (m *MockMetricsCollector) LogReadCount(arg0, arg1 string) prometheus.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogReadCount", arg0, arg1)
	ret0, _ := ret[0].(prometheus.Counter)
	return ret0
}

// LogReadCount indicates an expected call of LogReadCount
func (mr *MockMetricsCollectorMockRecorder) LogReadCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogReadCount", reflect.TypeOf((*MockMetricsCollector)(nil).LogReadCount), arg0, arg1)
}

// LogWriteCount mocks base method
func (m *MockMetricsCollector) LogWriteCount(arg0, arg1 string) prometheus.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogWriteCount", arg0, arg1)
	ret0, _ := ret[0].(prometheus.Counter)
	return ret0
}

// LogWriteCount indicates an expected call of LogWriteCount
func (mr *MockMetricsCollectorMockRecorder) LogWriteCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogWriteCount", reflect.TypeOf((*MockMetricsCollector)(nil).LogWriteCount), arg0, arg1)
}

// PingFailureCount mocks base method
func (m *MockMetricsCollector) PingFailureCount(arg0 string) prometheus.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingFailureCount", arg0)
	ret0, _ := ret[0].(prometheus.Counter)
	return ret0
}

// PingFailureCount indicates an expected call of PingFailureCount
func (mr *MockMetricsCollectorMockRecorder) PingFailureCount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingFailureCount", reflect.TypeOf((*MockMetricsCollector)(nil).PingFailureCount), arg0)
}

// TotalConnections mocks base method
func (m *MockMetricsCollector) TotalConnections() prometheus.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TotalConnections")
	ret0, _ := ret[0].(prometheus.Counter)
	return ret0
}

// TotalConnections indicates an expected call of TotalConnections
func (mr *MockMetricsCollectorMockRecorder) TotalConnections() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TotalConnections", reflect.TypeOf((*MockMetricsCollector)(nil).TotalConnections))
}
