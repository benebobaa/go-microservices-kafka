package provider

import (
	"context"
	"payment-svc/internal/dto"
	interfaces "payment-svc/internal/interface"
)

type PaymentProvider interface {
	ProcessPayment(ctx context.Context, req *dto.PaymentRequest) (*dto.BaseResponse[dto.Transaction], *dto.ErrorResponse)
	RefundPayment(ctx context.Context, req *dto.PaymentRequest) (*dto.BaseResponse[dto.Transaction], *dto.ErrorResponse)
	CreateAccountBalance(ctx context.Context, req *dto.AccountBalanceRequest) (*dto.BaseResponse[dto.AccountBalance], *dto.ErrorResponse)
}

type PaymentProviderImpl struct {
	client interfaces.Client
}

func NewPaymentProviderImpl(client interfaces.Client) PaymentProvider {
	return &PaymentProviderImpl{client: client}
}

func (u *PaymentProviderImpl) RefundPayment(ctx context.Context, req *dto.PaymentRequest) (*dto.BaseResponse[dto.Transaction], *dto.ErrorResponse) {
	var response dto.BaseResponse[dto.Transaction]

	//counter := 0
	//err := retryit.Do(ctx, func(ctx context.Context) error {
	//	counter++
	//	log.Println("retrying refund: ", counter)
	//	return u.client.PATCH(ctx, "/refund", req, &response)
	//}, retryit.WithInitialDelay(500*time.Millisecond))

	err := u.client.PATCH(ctx, "/refund", req, &response)

	if err != nil {
		return &response, &dto.ErrorResponse{Error: err.Error()}
	}

	if response.StatusCode != 200 {
		return &response, &dto.ErrorResponse{Error: response.Error}
	}

	return &response, nil
}

func (u *PaymentProviderImpl) CreateAccountBalance(ctx context.Context, req *dto.AccountBalanceRequest) (*dto.BaseResponse[dto.AccountBalance], *dto.ErrorResponse) {
	var response dto.BaseResponse[dto.AccountBalance]

	//counter := 0
	//err := retryit.Do(ctx, func(ctx context.Context) error {
	//	counter++
	//	log.Println("retrying create account balance: ", counter)
	//	return u.client.POST(ctx, "/balances", req, &response)
	//}, retryit.WithInitialDelay(500*time.Millisecond))

	err := u.client.POST(ctx, "/balances", req, &response)

	if err != nil {
		return &response, &dto.ErrorResponse{Error: err.Error()}
	}

	if response.StatusCode != 201 {
		return &response, &dto.ErrorResponse{Error: response.Error}
	}

	return &response, nil
}

func (u *PaymentProviderImpl) ProcessPayment(ctx context.Context, req *dto.PaymentRequest) (*dto.BaseResponse[dto.Transaction], *dto.ErrorResponse) {
	var response dto.BaseResponse[dto.Transaction]

	//counter := 0
	//err := retryit.Do(ctx, func(ctx context.Context) error {
	//	counter++
	//	log.Println("retrying payment: ", counter)
	//	return u.client.POST(ctx, "", req, &response)
	//}, retryit.WithInitialDelay(500*time.Millisecond))

	err := u.client.POST(ctx, "", req, &response)

	if err != nil {
		return &response, &dto.ErrorResponse{Error: err.Error()}
	}

	if response.StatusCode != 201 {
		return &response, &dto.ErrorResponse{Error: response.Error}
	}

	return &response, nil
}
