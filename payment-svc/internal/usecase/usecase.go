package usecase

import (
	"context"
	"log"
	"payment-svc/internal/dto"
	"payment-svc/internal/dto/event"
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

func (u *Usecase) ReserveProductMessaging(ctx context.Context, ge event.GlobalEvent[dto.ProductRequest]) error {

	return nil
}

func (u *Usecase) ProcessPayment(ctx context.Context, req *dto.ProductRequest) (*dto.ProductResponse, error) {
	var response dto.BaseResponse[dto.ProductResponse]
	log.Println("ReserveProduct: ", req)
	err := u.userClient.POST("/products/reserve", req, &response)

	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}
