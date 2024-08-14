package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

func (u *Usecase) ProcessPaymentMessaging(ctx context.Context, ge event.GlobalEvent[dto.PaymentRequest, any]) error {
	response, err := u.ProcessPayment(ctx, &dto.PaymentRequest{
		RefId:         ge.Payload.Request.RefId,
		Amount:        ge.Payload.Request.Amount,
		AccountBankID: ge.Payload.Request.AccountBankID,
	})

	var gevent event.GlobalEvent[dto.PaymentRequest, any]
	basePayload := event.BasePayload[dto.PaymentRequest, any]{
		Request: ge.Payload.Request,
	}

	if err != nil || response.Error != "" {
		if err != nil {
			basePayload.Response = err.Error()
		} else {
			basePayload.Response = response.Error
		}

		gevent = event.NewGlobalEvent[dto.PaymentRequest, any](
			"update",
			"error",
			"payment_failed",
			basePayload,
		)
	} else {
		basePayload.Response = response.Data

		gevent = event.NewGlobalEvent[dto.PaymentRequest, any](
			"update",
			"success",
			"payment_success",
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
		return fmt.Errorf("payment processing failed: %w", err)
	}

	return nil
}

func (u *Usecase) ProcessPayment(ctx context.Context, req *dto.PaymentRequest) (*dto.BaseResponse[dto.Transaction], error) {
	var response dto.BaseResponse[dto.Transaction]
	err := u.userClient.POST("/payment", req, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}
