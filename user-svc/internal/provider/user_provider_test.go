package provider_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"user-svc/internal/dto"
	"user-svc/internal/provider"
)

type mockUserClient struct {
	GetFunc   func(ctx context.Context, path string, req interface{}, resp interface{}) error
	PatchFunc func(ctx context.Context, path string, req interface{}, resp interface{}) error
	PostFunc  func(ctx context.Context, path string, req interface{}, resp interface{}) error
}

func (m *mockUserClient) GET(ctx context.Context, path string, req interface{}, resp interface{}) error {
	return m.GetFunc(ctx, path, req, resp)
}

func (m *mockUserClient) PATCH(ctx context.Context, path string, req interface{}, resp interface{}) error {
	return m.PatchFunc(ctx, path, req, resp)
}

func (m *mockUserClient) POST(ctx context.Context, path string, req interface{}, resp interface{}) error {
	return m.PostFunc(ctx, path, req, resp)
}

func TestUserProviderImpl_GetUserDetail(t *testing.T) {
	type fields struct {
		client *mockUserClient
	}
	type args struct {
		ctx context.Context
		req *dto.UserValidateRequest
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantResponse    *dto.BaseResponse[dto.UserResponse]
		wantErrResponse *dto.ErrorResponse
		wantErr         bool
	}{
		{
			name: "successful get user detail",
			fields: fields{
				client: &mockUserClient{
					GetFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.UserResponse])
						if !ok {
							return errors.New("unexpected response type")
						}
						r.StatusCode = 200
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.UserValidateRequest{Username: "testuser"},
			},
			wantResponse: &dto.BaseResponse[dto.UserResponse]{
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "failed get user detail",
			fields: fields{
				client: &mockUserClient{
					GetFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.UserResponse])
						if !ok {
							return errors.New("unexpected response type")
						}
						r.StatusCode = 400
						r.Error = "failed to get user detail"
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.UserValidateRequest{Username: "testuser"},
			},
			wantErrResponse: &dto.ErrorResponse{
				Message: "failed to get user detail",
			},
			wantErr: true,
		},
		{
			name: "client error",
			fields: fields{
				client: &mockUserClient{
					GetFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						return errors.New("client error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.UserValidateRequest{Username: "testuser"},
			},
			wantErrResponse: &dto.ErrorResponse{
				Message: "client error",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := provider.NewUserProvider(tt.fields.client)
			_, gotErrResponse := u.GetUserDetail(tt.args.ctx, tt.args.req)
			if (gotErrResponse != nil) != tt.wantErr {
				t.Errorf("UserProviderImpl.GetUserDetail() error = %v, wantErr %v", gotErrResponse, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
			//	t.Errorf("UserProviderImpl.GetUserDetail() = %v, want %v", gotResponse, tt.wantResponse)
			//}
			if !reflect.DeepEqual(gotErrResponse, tt.wantErrResponse) {
				t.Errorf("UserProviderImpl.GetUserDetail() error response = %v, want %v", gotErrResponse, tt.wantErrResponse)
			}
		})
	}
}

func TestUserProviderImpl_UpdateUser(t *testing.T) {
	type fields struct {
		client *mockUserClient
	}
	type args struct {
		ctx context.Context
		req *dto.UpdateBankIDRequest
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantResponse    *dto.BaseResponse[dto.UserResponse]
		wantErrResponse *dto.ErrorResponse
		wantErr         bool
	}{
		{
			name: "successful update user",
			fields: fields{
				client: &mockUserClient{
					PatchFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.UserResponse])
						if !ok {
							return errors.New("unexpected response type")
						}
						r.StatusCode = 200
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.UpdateBankIDRequest{Username: "testuser"},
			},
			wantResponse: &dto.BaseResponse[dto.UserResponse]{
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "failed update user",
			fields: fields{
				client: &mockUserClient{
					PatchFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.UserResponse])
						if !ok {
							return errors.New("unexpected response type")
						}
						r.StatusCode = 400
						r.Error = "failed to update user"
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.UpdateBankIDRequest{Username: "testuser"},
			},
			wantErrResponse: &dto.ErrorResponse{
				Message: "failed to update user",
			},
			wantErr: true,
		},
		{
			name: "client error",
			fields: fields{
				client: &mockUserClient{
					PatchFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						return errors.New("client error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.UpdateBankIDRequest{Username: "testuser"},
			},
			wantErrResponse: &dto.ErrorResponse{
				Message: "client error",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := provider.NewUserProvider(tt.fields.client)
			_, gotErrResponse := u.UpdateUser(tt.args.ctx, tt.args.req)
			if (gotErrResponse != nil) != tt.wantErr {
				t.Errorf("UserProviderImpl.UpdateUser() error = %v, wantErr %v", gotErrResponse, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
			//	t.Errorf("UserProviderImpl.UpdateUser() = %v, want %v", gotResponse, tt.wantResponse)
			//}
			if !reflect.DeepEqual(gotErrResponse, tt.wantErrResponse) {
				t.Errorf("UserProviderImpl.UpdateUser() error response = %v, want %v", gotErrResponse, tt.wantErrResponse)
			}
		})
	}
}

func TestUserProviderImpl_CreateUser(t *testing.T) {
	type fields struct {
		client *mockUserClient
	}
	type args struct {
		ctx context.Context
		req *dto.UserCreateRequest
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantResponse    *dto.BaseResponse[dto.UserResponse]
		wantErrResponse *dto.ErrorResponse
		wantErr         bool
	}{
		{
			name: "successful create user",
			fields: fields{
				client: &mockUserClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.UserResponse])
						if !ok {
							return errors.New("unexpected response type")
						}
						r.StatusCode = 201
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.UserCreateRequest{},
			},
			wantResponse: &dto.BaseResponse[dto.UserResponse]{
				StatusCode: 201,
			},
			wantErr: false,
		},
		{
			name: "failed create user",
			fields: fields{
				client: &mockUserClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.UserResponse])
						if !ok {
							return errors.New("unexpected response type")
						}
						r.StatusCode = 400
						r.Error = "failed to create user"
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.UserCreateRequest{},
			},
			wantErrResponse: &dto.ErrorResponse{
				Message: "failed to create user",
			},
			wantErr: true,
		},
		{
			name: "client error",
			fields: fields{
				client: &mockUserClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						return errors.New("client error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.UserCreateRequest{},
			},
			wantErrResponse: &dto.ErrorResponse{
				Message: "client error",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := provider.NewUserProvider(tt.fields.client)
			_, gotErrResponse := u.CreateUser(tt.args.ctx, tt.args.req)
			if (gotErrResponse != nil) != tt.wantErr {
				t.Errorf("UserProviderImpl.CreateUser() error = %v, wantErr %v", gotErrResponse, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
			//	t.Errorf("UserProviderImpl.CreateUser() = %v, want %v", gotResponse, tt.wantResponse)
			//}
			if !reflect.DeepEqual(gotErrResponse, tt.wantErrResponse) {
				t.Errorf("UserProviderImpl.CreateUser() error response = %v, want %v", gotErrResponse, tt.wantErrResponse)
			}
		})
	}
}
