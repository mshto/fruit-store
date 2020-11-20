// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mshto/fruit-store/bill (interfaces: Bill)

// Package billmock is a generated GoMock package.
package billmock

import (
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	bill "github.com/mshto/fruit-store/bill"
	config "github.com/mshto/fruit-store/config"
	entity "github.com/mshto/fruit-store/entity"
	reflect "reflect"
)

// MockBill is a mock of Bill interface
type MockBill struct {
	ctrl     *gomock.Controller
	recorder *MockBillMockRecorder
}

// MockBillMockRecorder is the mock recorder for MockBill
type MockBillMockRecorder struct {
	mock *MockBill
}

// NewMockBill creates a new mock instance
func NewMockBill(ctrl *gomock.Controller) *MockBill {
	mock := &MockBill{ctrl: ctrl}
	mock.recorder = &MockBillMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBill) EXPECT() *MockBillMockRecorder {
	return m.recorder
}

// GetDiscountByUser mocks base method
func (m *MockBill) GetDiscountByUser(arg0 uuid.UUID) (config.GeneralSale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDiscountByUser", arg0)
	ret0, _ := ret[0].(config.GeneralSale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDiscountByUser indicates an expected call of GetDiscountByUser
func (mr *MockBillMockRecorder) GetDiscountByUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDiscountByUser", reflect.TypeOf((*MockBill)(nil).GetDiscountByUser), arg0)
}

// GetTotalInfo mocks base method
func (m *MockBill) GetTotalInfo(arg0 uuid.UUID, arg1 []entity.GetUserProduct) (bill.TotalInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalInfo", arg0, arg1)
	ret0, _ := ret[0].(bill.TotalInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalInfo indicates an expected call of GetTotalInfo
func (mr *MockBillMockRecorder) GetTotalInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalInfo", reflect.TypeOf((*MockBill)(nil).GetTotalInfo), arg0, arg1)
}

// RemoveDiscount mocks base method
func (m *MockBill) RemoveDiscount(arg0 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveDiscount", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveDiscount indicates an expected call of RemoveDiscount
func (mr *MockBillMockRecorder) RemoveDiscount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDiscount", reflect.TypeOf((*MockBill)(nil).RemoveDiscount), arg0)
}

// SetDiscount mocks base method
func (m *MockBill) SetDiscount(arg0 uuid.UUID, arg1 config.GeneralSale) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDiscount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDiscount indicates an expected call of SetDiscount
func (mr *MockBillMockRecorder) SetDiscount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDiscount", reflect.TypeOf((*MockBill)(nil).SetDiscount), arg0, arg1)
}

// ValidateCard mocks base method
func (m *MockBill) ValidateCard(arg0 entity.Payment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateCard", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateCard indicates an expected call of ValidateCard
func (mr *MockBillMockRecorder) ValidateCard(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateCard", reflect.TypeOf((*MockBill)(nil).ValidateCard), arg0)
}
