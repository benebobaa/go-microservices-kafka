package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"payment-svc/internal/dto"
	"payment-svc/internal/dto/event"
	"payment-svc/internal/provider"
	"payment-svc/pkg/producer"
)

type Usecase struct {
	paymentProvider   provider.PaymentProvider
	orchestraProducer *producer.KafkaProducer
}

func NewUsecase(paymentProvider provider.PaymentProvider, orchestraProducer *producer.KafkaProducer) *Usecase {
	return &Usecase{
		paymentProvider:   paymentProvider,
		orchestraProducer: orchestraProducer,
	}
}

func (u *Usecase) CreateAccountBalanceMessaging(ctx context.Context, ge event.GlobalEvent[dto.AccountBalanceRequest, any]) error {
	response, err := u.paymentProvider.CreateAccountBalance(ctx, &dto.AccountBalanceRequest{
		Deposit:  ge.Payload.Request.Deposit,
		Username: ge.Payload.Request.Username,
	})

	var gevent event.GlobalEvent[dto.AccountBalanceRequest, any]

	basePayload := event.BasePayload[dto.AccountBalanceRequest, any]{
		Request: ge.Payload.Request,
	}

	if err != nil {
		basePayload.Response = err

		gevent = event.NewGlobalEvent[dto.AccountBalanceRequest, any](
			"create",
			"error",
			"bank_account_failed",
			basePayload,
		)

		if response.Error != "" {
			gevent.StatusCode = response.StatusCode
		} else {
			gevent.StatusCode = 500
		}

	} else {
		basePayload.Response = response.Data

		gevent = event.NewGlobalEvent[dto.AccountBalanceRequest, any](
			"create",
			"success",
			"bank_account_created",
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
		return fmt.Errorf("account balance processing failed: %s", err)
	}

	return nil
}

func (u *Usecase) ProcessPaymentMessaging(ctx context.Context, ge event.GlobalEvent[dto.PaymentRequest, any]) error {
	response, err := u.paymentProvider.ProcessPayment(ctx, &dto.PaymentRequest{
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
		return fmt.Errorf("payment processing failed: %s", err)
	}

	return nil
}

func (u *Usecase) RefundPaymentMessaging(ctx context.Context, ge event.GlobalEvent[dto.PaymentRequest, any]) error {
	response, err := u.paymentProvider.RefundPayment(ctx, &dto.PaymentRequest{
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
		return fmt.Errorf("refund processing failed: %s", err)
	}

	return nil
}
