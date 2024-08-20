package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"product-svc/internal/dto"
	"product-svc/internal/dto/event"
	"product-svc/internal/interfaces"
	"product-svc/internal/provider"
)

type Usecase struct {
	productProvider   provider.ProductProvider
	orchestraProducer interfaces.Producer
}

func NewUsecase(productProvider provider.ProductProvider, orchestraProducer interfaces.Producer) *Usecase {
	return &Usecase{
		productProvider:   productProvider,
		orchestraProducer: orchestraProducer,
	}
}

func (u *Usecase) ReserveProductMessaging(ctx context.Context, ge event.GlobalEvent[dto.ProductRequest, any]) error {
	response, err := u.productProvider.ReserveProduct(ctx, &dto.ProductRequest{
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
		return fmt.Errorf("product reservation failed: %s", err)
	}

	return nil
}

func (u *Usecase) ReleaseProductMessaging(ctx context.Context, ge event.GlobalEvent[dto.ProductRequest, any]) error {
	response, err := u.productProvider.ReleaseProduct(ctx, &dto.ProductRequest{
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
