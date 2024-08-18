package usecase

import (
	"context"
	"fmt"
	"log"
	"payment-svc/internal/dto"
	"payment-svc/internal/dto/event"
	"payment-svc/pkg/http_client"
	"payment-svc/pkg/producer"
	"time"

	"github.com/benebobaa/retry-it"
	"github.com/google/uuid"
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

	if err != nil {
		basePayload.Response = err

		gevent = event.NewGlobalEvent[dto.PaymentRequest, any](
			"update",
			"error",
			"payment_failed",
			basePayload,
		)
		if response.Error != "" {
			gevent.StatusCode = response.StatusCode
		} else {
			gevent.StatusCode = 500
		}
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

func (u *Usecase) RefundPaymentMessaging(ctx context.Context, ge event.GlobalEvent[dto.PaymentRequest, any]) error {
	response, err := u.RefundPayment(ctx, &dto.PaymentRequest{
		RefId: ge.Payload.Request.RefId,
	})

	var gevent event.GlobalEvent[dto.PaymentRequest, any]
	basePayload := event.BasePayload[dto.PaymentRequest, any]{
		Request: ge.Payload.Request,
	}

	if err != nil {
		basePayload.Response = err

		gevent = event.NewGlobalEvent[dto.PaymentRequest, any](
			"update",
			"error",
			"refund_failed",
			basePayload,
		)

		if response.Error != "" {
			gevent.StatusCode = response.StatusCode
		} else {
			gevent.StatusCode = 500
		}
	} else {
		basePayload.Response = response.Data

		gevent = event.NewGlobalEvent[dto.PaymentRequest, any](
			"update",
			"success",
			"refund_success",
			basePayload,
		)

		gevent.StatusCode = response.StatusCode
	}

	gevent.EventID = ge.EventID
	gevent.InstanceID = ge.InstanceID
	gevent.EventType = ge.EventType

	bytes, jsonErr := gevent.ToJSON()
	if jsonErr != nil {
		return fmt.Errorf("failed to convert event to JSON: %w", jsonErr)
	}

	sendErr := u.orchestraProducer.SendMessage(uuid.New().String(), bytes)
	if sendErr != nil {
		return fmt.Errorf("failed to send message: %w", sendErr)
	}

	if err != nil {
		return fmt.Errorf("refund processing failed: %w", err)
	}

	return nil
}

func (u *Usecase) ProcessPayment(ctx context.Context, req *dto.PaymentRequest) (*dto.BaseResponse[dto.Transaction], *dto.ErrorResponse) {
	var response dto.BaseResponse[dto.Transaction]

	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		log.Println("retrying payment: ", counter)
		return u.userClient.POST(ctx, "", req, &response)
	}, retryit.WithInitialDelay(500*time.Millisecond))

	if err != nil {
		return &response, &dto.ErrorResponse{Error: err.Error()}
	}

	if response.StatusCode != 201 {
		return &response, &dto.ErrorResponse{Error: response.Error}
	}

	return &response, nil
}

func (u *Usecase) RefundPayment(ctx context.Context, req *dto.PaymentRequest) (*dto.BaseResponse[dto.Transaction], *dto.ErrorResponse) {
	var response dto.BaseResponse[dto.Transaction]

	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		log.Println("retrying refund: ", counter)
		return u.userClient.PATCH(ctx, "/refund", req, &response)
	}, retryit.WithInitialDelay(500*time.Millisecond))

	if err != nil {
		return &response, &dto.ErrorResponse{Error: err.Error()}
	}

	if response.StatusCode != 200 {
		return &response, &dto.ErrorResponse{Error: response.Error}
	}

	return &response, nil
}
