package usecase

import (
	"context"
	"payment-svc/internal/dto"
	"payment-svc/pkg/http_client"
	"payment-svc/pkg/producer"
)

type Usecase struct {
	userClient        *http_client.PaymentClient
	orchestraProducer *producer.KafkaProducer
}

func NewUsecase(userClient *http_client.PaymentClient, orchestraProducer *producer.KafkaProducer) *Usecase {
	return &Usecase{
		userClient:        userClient,
		orchestraProducer: orchestraProducer,
	}
}

func (u *Usecase) ProcessPayment(ctx context.Context, req *dto.Payment) error {
	var response dto.BaseResponse[dto.Payment]

	err := u.userClient.POST("/payment", req, &response)

	if err != nil || response.Error != "" {
		return err
	}

	return nil
}
func (u *Usecase) RefundPayment(ctx context.Context, req *dto.RefundRequest) error {
	var response dto.BaseResponse[dto.Payment]

	err := u.userClient.POST("/refund", req, &response)

	if err != nil || response.Error != "" {
		return err
	}

	return nil
}
