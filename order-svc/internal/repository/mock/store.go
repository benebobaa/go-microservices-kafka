// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/sqlc/store.go
//
// Generated by this command:
//
//	mockgen -source=internal/repository/sqlc/store.go -package mockdb -destination=internal/repository/mock/store.go
//

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	sqlc "order-svc/internal/repository/sqlc"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CountByID mocks base method.
func (m *MockStore) CountByID(ctx context.Context, refID string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByID", ctx, refID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountByID indicates an expected call of CountByID.
func (mr *MockStoreMockRecorder) CountByID(ctx, refID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByID", reflect.TypeOf((*MockStore)(nil).CountByID), ctx, refID)
}

// CreateBankAccountRegistration mocks base method.
func (m *MockStore) CreateBankAccountRegistration(ctx context.Context, arg sqlc.CreateBankAccountRegistrationParams) (sqlc.BankAccountRegistration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBankAccountRegistration", ctx, arg)
	ret0, _ := ret[0].(sqlc.BankAccountRegistration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBankAccountRegistration indicates an expected call of CreateBankAccountRegistration.
func (mr *MockStoreMockRecorder) CreateBankAccountRegistration(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBankAccountRegistration", reflect.TypeOf((*MockStore)(nil).CreateBankAccountRegistration), ctx, arg)
}

// CreateOrder mocks base method.
func (m *MockStore) CreateOrder(ctx context.Context, arg sqlc.CreateOrderParams) (sqlc.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", ctx, arg)
	ret0, _ := ret[0].(sqlc.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockStoreMockRecorder) CreateOrder(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockStore)(nil).CreateOrder), ctx, arg)
}

// FindBankAccountRegistrationByUsernameOrEmail mocks base method.
func (m *MockStore) FindBankAccountRegistrationByUsernameOrEmail(ctx context.Context, arg sqlc.FindBankAccountRegistrationByUsernameOrEmailParams) (sqlc.BankAccountRegistration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindBankAccountRegistrationByUsernameOrEmail", ctx, arg)
	ret0, _ := ret[0].(sqlc.BankAccountRegistration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindBankAccountRegistrationByUsernameOrEmail indicates an expected call of FindBankAccountRegistrationByUsernameOrEmail.
func (mr *MockStoreMockRecorder) FindBankAccountRegistrationByUsernameOrEmail(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindBankAccountRegistrationByUsernameOrEmail", reflect.TypeOf((*MockStore)(nil).FindBankAccountRegistrationByUsernameOrEmail), ctx, arg)
}

// FindOrderByID mocks base method.
func (m *MockStore) FindOrderByID(ctx context.Context, id int32) (sqlc.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrderByID", ctx, id)
	ret0, _ := ret[0].(sqlc.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrderByID indicates an expected call of FindOrderByID.
func (mr *MockStoreMockRecorder) FindOrderByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrderByID", reflect.TypeOf((*MockStore)(nil).FindOrderByID), ctx, id)
}

// FindOrderByRefID mocks base method.
func (m *MockStore) FindOrderByRefID(ctx context.Context, refID string) (sqlc.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrderByRefID", ctx, refID)
	ret0, _ := ret[0].(sqlc.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrderByRefID indicates an expected call of FindOrderByRefID.
func (mr *MockStoreMockRecorder) FindOrderByRefID(ctx, refID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrderByRefID", reflect.TypeOf((*MockStore)(nil).FindOrderByRefID), ctx, refID)
}

// FindOrdersByUsername mocks base method.
func (m *MockStore) FindOrdersByUsername(ctx context.Context, username string) ([]sqlc.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrdersByUsername", ctx, username)
	ret0, _ := ret[0].([]sqlc.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrdersByUsername indicates an expected call of FindOrdersByUsername.
func (mr *MockStoreMockRecorder) FindOrdersByUsername(ctx, username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrdersByUsername", reflect.TypeOf((*MockStore)(nil).FindOrdersByUsername), ctx, username)
}

// UpdateBankAccountRegistration mocks base method.
func (m *MockStore) UpdateBankAccountRegistration(ctx context.Context, arg sqlc.UpdateBankAccountRegistrationParams) (sqlc.BankAccountRegistration, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBankAccountRegistration", ctx, arg)
	ret0, _ := ret[0].(sqlc.BankAccountRegistration)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBankAccountRegistration indicates an expected call of UpdateBankAccountRegistration.
func (mr *MockStoreMockRecorder) UpdateBankAccountRegistration(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBankAccountRegistration", reflect.TypeOf((*MockStore)(nil).UpdateBankAccountRegistration), ctx, arg)
}

// UpdateOrder mocks base method.
func (m *MockStore) UpdateOrder(ctx context.Context, arg sqlc.UpdateOrderParams) (sqlc.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", ctx, arg)
	ret0, _ := ret[0].(sqlc.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockStoreMockRecorder) UpdateOrder(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockStore)(nil).UpdateOrder), ctx, arg)
}
