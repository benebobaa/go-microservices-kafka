package provider_test

import (
	"context"
	"errors"
	"payment-svc/internal/dto"
	"payment-svc/internal/provider"
	"reflect"
	"testing"
)

type mockPaymentClient struct {
	PatchFunc func(ctx context.Context, path string, req interface{}, resp interface{}) error
	PostFunc  func(ctx context.Context, path string, req interface{}, resp interface{}) error
	GetFunc   func(ctx context.Context, url string, request interface{}, response interface{}) error
}

func (m *mockPaymentClient) GET(ctx context.Context, url string, request interface{}, response interface{}) error {
	return m.GetFunc(ctx, url, request, response)
}

func (m *mockPaymentClient) PATCH(ctx context.Context, path string, req interface{}, resp interface{}) error {
	return m.PatchFunc(ctx, path, req, resp)
}

func (m *mockPaymentClient) POST(ctx context.Context, path string, req interface{}, resp interface{}) error {
	return m.PostFunc(ctx, path, req, resp)
}

func TestPaymentProviderImpl_RefundPayment(t *testing.T) {
	type fields struct {
		client *mockPaymentClient
	}
	type args struct {
		ctx context.Context
		req *dto.PaymentRequest
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantResponse    *dto.BaseResponse[dto.Transaction]
		wantErrResponse *dto.ErrorResponse
		wantErr         bool
	}{
		{
			name: "successful refund payment",
			fields: fields{
				client: &mockPaymentClient{
					PatchFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.Transaction])
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
				req: &dto.PaymentRequest{Amount: 100},
			},
			wantResponse: &dto.BaseResponse[dto.Transaction]{
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "failed refund payment",
			fields: fields{
				client: &mockPaymentClient{
					PatchFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.Transaction])
						if !ok {
							return errors.New("unexpected response type")
						}
						r.StatusCode = 400
						r.Error = "failed to refund payment"
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.PaymentRequest{Amount: 100},
			},
			wantErrResponse: &dto.ErrorResponse{
				Error: "failed to refund payment",
			},
			wantErr: true,
		},
		{
			name: "client error during refund",
			fields: fields{
				client: &mockPaymentClient{
					PatchFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						return errors.New("client error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.PaymentRequest{Amount: 100},
			},
			wantErrResponse: &dto.ErrorResponse{
				Error: "client error",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := provider.NewPaymentProviderImpl(tt.fields.client)
			_, gotErrResponse := p.RefundPayment(tt.args.ctx, tt.args.req)
			if (gotErrResponse != nil) != tt.wantErr {
				t.Errorf("PaymentProviderImpl.RefundPayment() error = %v, wantErr %v", gotErrResponse, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
			//	t.Errorf("PaymentProviderImpl.RefundPayment() = %v, want %v", gotResponse, tt.wantResponse)
			//}
			if !reflect.DeepEqual(gotErrResponse, tt.wantErrResponse) {
				t.Errorf("PaymentProviderImpl.RefundPayment() error response = %v, want %v", gotErrResponse, tt.wantErrResponse)
			}
		})
	}
}

func TestPaymentProviderImpl_CreateAccountBalance(t *testing.T) {
	type fields struct {
		client *mockPaymentClient
	}
	type args struct {
		ctx context.Context
		req *dto.AccountBalanceRequest
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantResponse    *dto.BaseResponse[dto.AccountBalance]
		wantErrResponse *dto.ErrorResponse
		wantErr         bool
	}{
		{
			name: "successful create account balance",
			fields: fields{
				client: &mockPaymentClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.AccountBalance])
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
				req: &dto.AccountBalanceRequest{
					Username: "bene",
					Deposit:  1000,
				},
			},
			wantResponse: &dto.BaseResponse[dto.AccountBalance]{
				StatusCode: 201,
			},
			wantErr: false,
		},
		{
			name: "failed create account balance",
			fields: fields{
				client: &mockPaymentClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.AccountBalance])
						if !ok {
							return errors.New("unexpected response type")
						}
						r.StatusCode = 400
						r.Error = "failed to create account balance"
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.AccountBalanceRequest{
					Username: "bene",
					Deposit:  1000,
				},
			},
			wantErrResponse: &dto.ErrorResponse{
				Error: "failed to create account balance",
			},
			wantErr: true,
		},
		{
			name: "client error during create account balance",
			fields: fields{
				client: &mockPaymentClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						return errors.New("client error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.AccountBalanceRequest{
					Username: "bene",
					Deposit:  1000,
				},
			},
			wantErrResponse: &dto.ErrorResponse{
				Error: "client error",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := provider.NewPaymentProviderImpl(tt.fields.client)
			_, gotErrResponse := p.CreateAccountBalance(tt.args.ctx, tt.args.req)
			if (gotErrResponse != nil) != tt.wantErr {
				t.Errorf("PaymentProviderImpl.CreateAccountBalance() error = %v, wantErr %v", gotErrResponse, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
			//	t.Errorf("PaymentProviderImpl.CreateAccountBalance() = %v, want %v", gotResponse, tt.wantResponse)
			//}
			if !reflect.DeepEqual(gotErrResponse, tt.wantErrResponse) {
				t.Errorf("PaymentProviderImpl.CreateAccountBalance() error response = %v, want %v", gotErrResponse, tt.wantErrResponse)
			}
		})
	}
}

func TestPaymentProviderImpl_ProcessPayment(t *testing.T) {
	type fields struct {
		client *mockPaymentClient
	}
	type args struct {
		ctx context.Context
		req *dto.PaymentRequest
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantResponse    *dto.BaseResponse[dto.Transaction]
		wantErrResponse *dto.ErrorResponse
		wantErr         bool
	}{
		{
			name: "successful process payment",
			fields: fields{
				client: &mockPaymentClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.Transaction])
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
				req: &dto.PaymentRequest{Amount: 100},
			},
			wantResponse: &dto.BaseResponse[dto.Transaction]{
				StatusCode: 201,
			},
			wantErr: false,
		},
		{
			name: "failed process payment",
			fields: fields{
				client: &mockPaymentClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						r, ok := resp.(*dto.BaseResponse[dto.Transaction])
						if !ok {
							return errors.New("unexpected response type")
						}
						r.StatusCode = 400
						r.Error = "failed to process payment"
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.PaymentRequest{Amount: 100},
			},
			wantErrResponse: &dto.ErrorResponse{
				Error: "failed to process payment",
			},
			wantErr: true,
		},
		{
			name: "client error during process payment",
			fields: fields{
				client: &mockPaymentClient{
					PostFunc: func(ctx context.Context, path string, req interface{}, resp interface{}) error {
						return errors.New("client error")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &dto.PaymentRequest{Amount: 100},
			},
			wantErrResponse: &dto.ErrorResponse{
				Error: "client error",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := provider.NewPaymentProviderImpl(tt.fields.client)
			_, gotErrResponse := p.ProcessPayment(tt.args.ctx, tt.args.req)
			if (gotErrResponse != nil) != tt.wantErr {
				t.Errorf("PaymentProviderImpl.ProcessPayment() error = %v, wantErr %v", gotErrResponse, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
			//	t.Errorf("PaymentProviderImpl.ProcessPayment() = %v, want %v", gotResponse, tt.wantResponse)
			//}
			if !reflect.DeepEqual(gotErrResponse, tt.wantErrResponse) {
				t.Errorf("PaymentProviderImpl.ProcessPayment() error response = %v, want %v", gotErrResponse, tt.wantErrResponse)
			}
		})
	}
}
