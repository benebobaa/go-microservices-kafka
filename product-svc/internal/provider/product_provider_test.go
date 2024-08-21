package provider_test

import (
	"context"
	"errors"
	"product-svc/internal/dto"
	"product-svc/internal/provider"
	"reflect"
	"testing"
)

type mockProductClient struct {
	PostFunc func(ctx context.Context, path string, req interface{}, resp interface{}) error
}

func (m *mockProductClient) POST(ctx context.Context, path string, req interface{}, resp interface{}) error {
	return m.PostFunc(ctx, path, req, resp)
}
func TestProductProviderImpl_ReserveProduct(t *testing.T) {
	type fields struct {
		client *mockProductClient
	}
	type args struct {
		ctx context.Context
		req *dto.ProductRequest
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantResponse    *dto.BaseResponse[dto.ProductResponse]
		wantErrResponse *dto.ErrorResponse
		wantErr         bool
	}{
		{
			name: "successful reserve",
			fields: fields{
				client: &mockProductClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.ProductResponse])
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
				req: &dto.ProductRequest{},
			},
			wantResponse: &dto.BaseResponse[dto.ProductResponse]{
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "failed reserve",
			fields: fields{
				client: &mockProductClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.ProductResponse])
						if !ok {
							return errors.New("unexpected response type")
						}
						r.StatusCode = 400
						r.Error = "failed to reserve product"
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.ProductRequest{},
			},
			wantErrResponse: &dto.ErrorResponse{
				Error: "failed to reserve product",
			},
			wantErr: true, // This should be true as you're expecting an error
		},
		{
			name: "client error",
			fields: fields{
				client: &mockProductClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						return errors.New("client error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.ProductRequest{},
			},
			wantErrResponse: &dto.ErrorResponse{
				Error: "client error",
			},
			wantErr: true, // This should be true as you're expecting an error
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := provider.NewProductProviderImpl(tt.fields.client)
			_, gotErrResponse := u.ReserveProduct(tt.args.ctx, tt.args.req)
			if (gotErrResponse != nil) != tt.wantErr {
				t.Errorf("ProductProviderImpl.ReserveProduct() error = %v, wantErr %v", gotErrResponse, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
			//	t.Errorf("ProductProviderImpl.ReserveProduct() = %v, want %v", gotResponse, tt.wantResponse)
			//}
			if !reflect.DeepEqual(gotErrResponse, tt.wantErrResponse) {
				t.Errorf("ProductProviderImpl.ReserveProduct() error response = %v, want %v", gotErrResponse, tt.wantErrResponse)
			}
		})
	}
}

func TestProductProviderImpl_ReleaseProduct(t *testing.T) {
	type fields struct {
		client *mockProductClient
	}
	type args struct {
		ctx context.Context
		req *dto.ProductRequest
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantResponse    *dto.BaseResponse[dto.ProductResponse]
		wantErrResponse *dto.ErrorResponse
		wantErr         bool
	}{
		{
			name: "successful release",
			fields: fields{
				client: &mockProductClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.ProductResponse])
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
				req: &dto.ProductRequest{},
			},
			wantResponse: &dto.BaseResponse[dto.ProductResponse]{
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "failed release",
			fields: fields{
				client: &mockProductClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.ProductResponse])
						if !ok {
							return errors.New("unexpected response type")
						}
						r.StatusCode = 400
						r.Error = "failed to release product"
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.ProductRequest{},
			},
			wantErrResponse: &dto.ErrorResponse{
				Error: "failed to release product",
			},
			wantErr: true, // This should be true as you're expecting an error
		},
		{
			name: "client error",
			fields: fields{
				client: &mockProductClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						return errors.New("client error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.ProductRequest{},
			},
			wantErrResponse: &dto.ErrorResponse{
				Error: "client error",
			},
			wantErr: true, // This should be true as you're expecting an error
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := provider.NewProductProviderImpl(tt.fields.client)
			_, gotErrResponse := u.ReleaseProduct(tt.args.ctx, tt.args.req)
			if (gotErrResponse != nil) != tt.wantErr {
				t.Errorf("ProductProviderImpl.ReleaseProduct() error = %v, wantErr %v", gotErrResponse, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
			//	t.Errorf("ProductProviderImpl.ReleaseProduct() = %v, want %v", gotResponse, tt.wantResponse)
			//}
			//if !reflect.DeepEqual(gotErrResponse, tt.wantErrResponse) {
			//	t.Errorf("ProductProviderImpl.ReleaseProduct() error response = %v, want %v", gotErrResponse, tt.wantErrResponse)
			//}
		})
	}
}
