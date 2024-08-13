package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"product-svc/internal/dto"
	"product-svc/internal/dto/event"
	"product-svc/pkg/http_client"
	"product-svc/pkg/producer"
)

type Usecase struct {
	userClient        *http_client.ProductClient
	orchestraProducer *producer.KafkaProducer
}

func NewUsecase(userClient *http_client.ProductClient, orchestraProducer *producer.KafkaProducer) *Usecase {
	return &Usecase{
		userClient:        userClient,
		orchestraProducer: orchestraProducer,
	}
}

func (u *Usecase) ReserveProductMessaging(ctx context.Context, ge event.GlobalEvent[dto.ProductRequest]) error {

	log.Println("ReserveProductMessaging: ", ge.Payload)

	response, err := u.ReserveProduct(ctx, &dto.ProductRequest{
		ProductID: ge.Payload.ProductID,
		Quantity:  ge.Payload.Quantity,
	})

	var gevent event.GlobalEvent[any]
	if err != nil || response.Error != "" {
		gevent = event.NewGlobalEvent[any](
			"update",
			"error",
			"product_reservation_failed",
			"product_reservation_error",
			response.Error, response.StatusCode,
		)
	} else {
		gevent = event.NewGlobalEvent[any](
			"update",
			"success",
			"product_reserved",
			"product_reservation_success",
			response.Data, response.StatusCode,
		)
	}

	gevent.EventID = ge.EventID
	gevent.InstanceID = ge.InstanceID
	gevent.EventType = ge.EventType
	gevent.StatusCode = response.StatusCode

	bytes, jsonErr := gevent.ToJSON()
	if jsonErr != nil {
		return fmt.Errorf("failed to convert event to JSON: %w", jsonErr)
	}

	sendErr := u.orchestraProducer.SendMessage(uuid.New().String(), bytes)
	if sendErr != nil {
		return fmt.Errorf("failed to send message: %w", sendErr)
	}

	if err != nil {
		return fmt.Errorf("product reservation failed: %w", err)
	}

	return nil
}

func (u *Usecase) ReserveProduct(ctx context.Context, req *dto.ProductRequest) (*dto.BaseResponse[dto.ProductResponse], error) {
	var response dto.BaseResponse[dto.ProductResponse]
	log.Println("ReserveProduct: ", req)
	err := u.userClient.POST("/products/reserve", req, &response)

	log.Println("ReserveProduct response: ", response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
