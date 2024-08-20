package usecase

import (
	"context"
	"github.com/google/uuid"
	"order-svc/internal/dto"
	"order-svc/internal/dto/event"
	"order-svc/internal/repository/sqlc"
	"order-svc/pkg/producer"
)

type BankRegistrationUsecase struct {
	queries           sqlc.Store
	orchestraProducer *producer.KafkaProducer
}

func NewBankRegistrationUsecase(queries sqlc.Store, producer *producer.KafkaProducer) *BankRegistrationUsecase {
	return &BankRegistrationUsecase{
		queries:           queries,
		orchestraProducer: producer,
	}
}

func (b *BankRegistrationUsecase) RegisterBankAccount(ctx context.Context, req *dto.BankRegistrationRequest) (*sqlc.BankAccountRegistration, error) {

	bankAccount, err := b.queries.CreateBankAccountRegistration(ctx, sqlc.CreateBankAccountRegistrationParams{
		CustomerID: uuid.New().String(),
		Username:   req.Username,
		Email:      req.Email,
		Status:     dto.PROCESSING.String(),
		Deposit:    req.Deposit,
	})

	if err != nil {
		return nil, err
	}

	basePayload := event.BasePayload[dto.BankRegistrationRequest, sqlc.BankAccountRegistration]{
		Request:  *req,
		Response: bankAccount,
	}

	orderEvent := event.NewGlobalEvent(
		"create",
		"success",
		"bank_regis_created",
		event.BANK_ACCOUNT_REGISTRATION.String(),
		basePayload,
	)
	orderEvent.StatusCode = 201
	bytes, err := orderEvent.ToJSON()

	if err != nil {
		return nil, err
	}

	err = b.orchestraProducer.SendMessage(uuid.New().String(), bytes)

	if err != nil {
		return nil, err
	}

	return &bankAccount, nil
}

func (b *BankRegistrationUsecase) UpdateBankRegistrationMessaging(ctx context.Context, req event.GlobalEvent[dto.BankRegistrationUpdate, any]) error {

	updatedRegis, err := b.queries.UpdateBankAccountRegistration(ctx, sqlc.UpdateBankAccountRegistrationParams{
		CustomerID: req.Payload.Request.CustomerID,
		Username:   req.Payload.Request.Username,
		Email:      req.Payload.Request.Email,
		Status:     req.Payload.Request.Status,
	})

	if err != nil {
		return err
	}

	basePayload := event.BasePayload[dto.BankRegistrationUpdate, any]{
		Request:  req.Payload.Request,
		Response: updatedRegis,
	}

	registEvent := event.NewGlobalEvent(
		"update",
		"success",
		"bank_regis_updated",
		event.BANK_ACCOUNT_REGISTRATION.String(),
		basePayload,
	)

	registEvent.EventID = req.EventID
	registEvent.InstanceID = req.InstanceID
	registEvent.EventType = req.EventType

	registEvent.StatusCode = 200
	bytes, err := registEvent.ToJSON()

	if err != nil {
		return err
	}

	err = b.orchestraProducer.SendMessage(uuid.New().String(), bytes)

	if err != nil {
		return err
	}

	return nil
}
