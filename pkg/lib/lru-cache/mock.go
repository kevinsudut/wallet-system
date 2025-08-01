// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/lib/lru-cache/interfaces.go
//
// Generated by this command:
//
//	mockgen -source=pkg/lib/lru-cache/interfaces.go -destination=pkg/lib/lru-cache/mock.go -package=lrucache
//

// Package lrucache is a generated GoMock package.
package lrucache

import (
	reflect "reflect"
	time "time"

	ccache "github.com/karlseguin/ccache/v3"
	gomock "go.uber.org/mock/gomock"
)

// MockLRUCacheItf is a mock of LRUCacheItf interface.
type MockLRUCacheItf struct {
	ctrl     *gomock.Controller
	recorder *MockLRUCacheItfMockRecorder
}

// MockLRUCacheItfMockRecorder is the mock recorder for MockLRUCacheItf.
type MockLRUCacheItfMockRecorder struct {
	mock *MockLRUCacheItf
}

// NewMockLRUCacheItf creates a new mock instance.
func NewMockLRUCacheItf(ctrl *gomock.Controller) *MockLRUCacheItf {
	mock := &MockLRUCacheItf{ctrl: ctrl}
	mock.recorder = &MockLRUCacheItfMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLRUCacheItf) EXPECT() *MockLRUCacheItfMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockLRUCacheItf) Delete(key string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", key)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockLRUCacheItfMockRecorder) Delete(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockLRUCacheItf)(nil).Delete), key)
}

// Fetch mocks base method.
func (m *MockLRUCacheItf) Fetch(key string, duration time.Duration, fetch func() (any, error)) (*ccache.Item[any], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", key, duration, fetch)
	ret0, _ := ret[0].(*ccache.Item[any])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockLRUCacheItfMockRecorder) Fetch(key, duration, fetch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockLRUCacheItf)(nil).Fetch), key, duration, fetch)
}

// Get mocks base method.
func (m *MockLRUCacheItf) Get(key string) *ccache.Item[any] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(*ccache.Item[any])
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockLRUCacheItfMockRecorder) Get(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockLRUCacheItf)(nil).Get), key)
}

// Set mocks base method.
func (m *MockLRUCacheItf) Set(key string, value any, duration time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", key, value, duration)
}

// Set indicates an expected call of Set.
func (mr *MockLRUCacheItfMockRecorder) Set(key, value, duration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockLRUCacheItf)(nil).Set), key, value, duration)
}
