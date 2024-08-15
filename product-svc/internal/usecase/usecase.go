package usecase

import (
	"context"
	"fmt"
	"log"
	"product-svc/internal/dto"
	"product-svc/internal/dto/event"
	"product-svc/pkg/http_client"
	"product-svc/pkg/producer"

	"github.com/google/uuid"
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

func (u *Usecase) ReserveProductMessaging(ctx context.Context, ge event.GlobalEvent[dto.ProductRequest, any]) error {
	response, err := u.ReserveProduct(ctx, &dto.ProductRequest{
		ProductID: ge.Payload.Request.ProductID,
		Quantity:  ge.Payload.Request.Quantity,
	})

	var gevent event.GlobalEvent[dto.ProductRequest, any]
	basePayload := event.BasePayload[dto.ProductRequest, any]{
		Request: ge.Payload.Request,
	}
	if err != nil || response.Error != "" {
		if err != nil {
			basePayload.Response = err.Error()
		} else {
			basePayload.Response = response.Error
		}

		gevent = event.NewGlobalEvent[dto.ProductRequest, any](
			"update",
			"error",
			"product_reservation_failed",
			basePayload,
		)
	} else {
		basePayload.Response = response.Data

		gevent = event.NewGlobalEvent[dto.ProductRequest, any](
			"update",
			"success",
			"product_reservation_success",
			basePayload,
		)
	}

	log.Println("check state:", gevent.State)

	gevent.EventID = ge.EventID
	gevent.InstanceID = ge.InstanceID
	gevent.EventType = ge.EventType
	if response != nil {
		gevent.StatusCode = response.StatusCode
	} else {
		gevent.StatusCode = 500
	}

	log.Println("gevent:", gevent)

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
	err := u.userClient.POST("/reserve", req, &response)

	log.Println("ReserveProduct response: ", response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
