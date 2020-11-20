// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mshto/fruit-store/authentication (interfaces: Auth)

// Package authmock is a generated GoMock package.
package authmock

import (
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	authentication "github.com/mshto/fruit-store/authentication"
	entity "github.com/mshto/fruit-store/entity"
	reflect "reflect"
)

// MockAuth is a mock of Auth interface
type MockAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAuthMockRecorder
}

// MockAuthMockRecorder is the mock recorder for MockAuth
type MockAuthMockRecorder struct {
	mock *MockAuth
}

// NewMockAuth creates a new mock instance
func NewMockAuth(ctrl *gomock.Controller) *MockAuth {
	mock := &MockAuth{ctrl: ctrl}
	mock.recorder = &MockAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuth) EXPECT() *MockAuthMockRecorder {
	return m.recorder
}

// CreateTokens mocks base method
func (m *MockAuth) CreateTokens(arg0 uuid.UUID) (*entity.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTokens", arg0)
	ret0, _ := ret[0].(*entity.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTokens indicates an expected call of CreateTokens
func (mr *MockAuthMockRecorder) CreateTokens(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTokens", reflect.TypeOf((*MockAuth)(nil).CreateTokens), arg0)
}

// GetUserUUID mocks base method
func (m *MockAuth) GetUserUUID(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserUUID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserUUID indicates an expected call of GetUserUUID
func (mr *MockAuthMockRecorder) GetUserUUID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserUUID", reflect.TypeOf((*MockAuth)(nil).GetUserUUID), arg0)
}

// RefreshTokens mocks base method
func (m *MockAuth) RefreshTokens(arg0 string) (*entity.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshTokens", arg0)
	ret0, _ := ret[0].(*entity.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshTokens indicates an expected call of RefreshTokens
func (mr *MockAuthMockRecorder) RefreshTokens(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshTokens", reflect.TypeOf((*MockAuth)(nil).RefreshTokens), arg0)
}

// RemoveTokens mocks base method
func (m *MockAuth) RemoveTokens(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveTokens", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveTokens indicates an expected call of RemoveTokens
func (mr *MockAuthMockRecorder) RemoveTokens(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveTokens", reflect.TypeOf((*MockAuth)(nil).RemoveTokens), arg0, arg1)
}

// ValidateToken mocks base method
func (m *MockAuth) ValidateToken(arg0 string) (*authentication.AccessDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", arg0)
	ret0, _ := ret[0].(*authentication.AccessDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateToken indicates an expected call of ValidateToken
func (mr *MockAuthMockRecorder) ValidateToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockAuth)(nil).ValidateToken), arg0)
}
