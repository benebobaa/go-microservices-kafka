package usecase

import (
	"context"
	"fmt"
	"log"
	"product-svc/internal/dto"
	"product-svc/internal/dto/event"
	"product-svc/pkg/http_client"
	"product-svc/pkg/producer"
	"time"

	"github.com/benebobaa/retry-it"
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
			basePayload.Response = err
		} else {
			basePayload.Response = response
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

func (u *Usecase) ReleaseProductMessaging(ctx context.Context, ge event.GlobalEvent[dto.ProductRequest, any]) error {
	response, err := u.ReleaseProduct(ctx, &dto.ProductRequest{
		ProductID: ge.Payload.Request.ProductID,
		Quantity:  ge.Payload.Request.Quantity,
	})

	var gevent event.GlobalEvent[dto.ProductRequest, any]
	basePayload := event.BasePayload[dto.ProductRequest, any]{
		Request: ge.Payload.Request,
	}

	if err != nil {
		basePayload.Response = err

		gevent = event.NewGlobalEvent[dto.ProductRequest, any](
			"update",
			"error",
			"product_release_failed",
			basePayload,
		)

		if response.Error != "" {
			gevent.StatusCode = response.StatusCode
		} else {
			gevent.StatusCode = 500
		}
	} else {
		basePayload.Response = response.Data

		gevent = event.NewGlobalEvent[dto.ProductRequest, any](
			"update",
			"success",
			"product_release_success",
			basePayload,
		)

		gevent.StatusCode = response.StatusCode
	}

	log.Println("check state:", gevent.State)

	gevent.EventID = ge.EventID
	gevent.InstanceID = ge.InstanceID
	gevent.EventType = ge.EventType

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
		return fmt.Errorf("product release failed: %s", err)
	}

	return nil
}

func (u *Usecase) ReserveProduct(ctx context.Context, req *dto.ProductRequest) (*dto.BaseResponse[dto.ProductResponse], *dto.ErrorResponse) {
	var response dto.BaseResponse[dto.ProductResponse]

	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		log.Println("retrying reserve: ", counter)
		return u.userClient.POST(ctx, "/reserve", req, &response)
	}, retryit.WithInitialDelay(500*time.Millisecond))

	if err != nil {
		return &response, &dto.ErrorResponse{Error: err.Error()}
	}

	if response.StatusCode != 200 {
		return &response, &dto.ErrorResponse{Error: response.Error}
	}

	return &response, nil
}

func (u *Usecase) ReleaseProduct(ctx context.Context, req *dto.ProductRequest) (*dto.BaseResponse[dto.ProductResponse], *dto.ErrorResponse) {
	var response dto.BaseResponse[dto.ProductResponse]

	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		log.Println("retrying release: ", counter)
		return u.userClient.POST(ctx, "/release", req, &response)
	}, retryit.WithInitialDelay(500*time.Millisecond))

	if err != nil {
		return &response, &dto.ErrorResponse{Error: err.Error()}
	}

	if response.StatusCode != 200 {
		return &response, &dto.ErrorResponse{Error: response.Error}
	}

	return &response, nil
}
