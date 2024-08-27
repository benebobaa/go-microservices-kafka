package provider

import (
	"context"
	"github.com/golang/mock/gomock"
	"user-svc/internal/dto"
)

type MockUserProvider struct {
	ctrl     *gomock.Controller
	recorder *MockUserProviderMockRecorder
}

func (m *MockUserProvider) SendMessage(key string, message []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", key, message)
	ret0, _ := ret[0].(error)
	return ret0
}

type MockUserProviderMockRecorder struct {
	mock *MockUserProvider
}

func NewMockUserProvider(ctrl *gomock.Controller) *MockUserProvider {
	mock := &MockUserProvider{ctrl: ctrl}
	mock.recorder = &MockUserProviderMockRecorder{mock}
	return mock
}

func (m *MockUserProvider) EXPECT() *MockUserProviderMockRecorder {
	return m.recorder
}

// GetUserDetail mocks base method.
func (m *MockUserProvider) GetUserDetail(ctx context.Context, request *dto.UserValidateRequest) (*dto.BaseResponse[dto.UserResponse], *dto.ErrorResponse) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDetail", ctx, request)
	ret0, _ := ret[0].(*dto.BaseResponse[dto.UserResponse])
	ret1, _ := ret[1].(*dto.ErrorResponse)
	return ret0, ret1
}

// UpdateUser mocks base method.
func (m *MockUserProvider) UpdateUser(ctx context.Context, request *dto.UpdateBankIDRequest) (*dto.BaseResponse[dto.UserResponse], *dto.ErrorResponse) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, request)
	ret0, _ := ret[0].(*dto.BaseResponse[dto.UserResponse])
	ret1, _ := ret[1].(*dto.ErrorResponse)
	return ret0, ret1
}

// CreateUser mocks base method.
func (m *MockUserProvider) CreateUser(ctx context.Context, request *dto.UserCreateRequest) (*dto.BaseResponse[dto.UserResponse], *dto.ErrorResponse) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, request)
	ret0, _ := ret[0].(*dto.BaseResponse[dto.UserResponse])
	ret1, _ := ret[1].(*dto.ErrorResponse)
	return ret0, ret1
}
