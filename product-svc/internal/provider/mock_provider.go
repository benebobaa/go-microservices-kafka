package provider

import (
	"context"
	"github.com/golang/mock/gomock"
	"product-svc/internal/dto"
	"reflect"
)

type MockProductProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProductProviderMockRecorder
}

type MockProductProviderMockRecorder struct {
	mock *MockProductProvider
}

func NewMockProductProvider(ctrl *gomock.Controller) *MockProductProvider {
	mock := &MockProductProvider{ctrl: ctrl}
	mock.recorder = &MockProductProviderMockRecorder{mock}
	return mock
}

func (m *MockProductProvider) EXPECT() *MockProductProviderMockRecorder {
	return m.recorder
}

func (m *MockProductProvider) ReserveProduct(ctx context.Context, req *dto.ProductRequest) (*dto.BaseResponse[dto.ProductResponse], *dto.ErrorResponse) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReserveProduct", ctx, req)
	ret0, _ := ret[0].(*dto.BaseResponse[dto.ProductResponse])
	ret1, _ := ret[1].(*dto.ErrorResponse)
	return ret0, ret1
}

func (mr *MockProductProviderMockRecorder) ReserveProduct(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReserveProduct", reflect.TypeOf((*MockProductProvider)(nil).ReserveProduct), ctx, req)
}

func (m *MockProductProvider) ReleaseProduct(ctx context.Context, req *dto.ProductRequest) (*dto.BaseResponse[dto.ProductResponse], *dto.ErrorResponse) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReleaseProduct", ctx, req)
	ret0, _ := ret[0].(*dto.BaseResponse[dto.ProductResponse])
	ret1, _ := ret[1].(*dto.ErrorResponse)
	return ret0, ret1
}

func (mr *MockProductProviderMockRecorder) ReleaseProduct(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseProduct", reflect.TypeOf((*MockProductProvider)(nil).ReleaseProduct), ctx, req)
}
