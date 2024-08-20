package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"user-svc/internal/dto"
	"user-svc/internal/dto/event"
	"user-svc/internal/interfaces"
	"user-svc/internal/provider"
)

type Usecase struct {
	userClient        interfaces.Client
	orchestraProducer interfaces.Producer
	provider          provider.UserProvider
}

func NewUsecase(
	userClient interfaces.Client,
	orchestraProducer interfaces.Producer,
	provider provider.UserProvider,
) *Usecase {
	return &Usecase{
		userClient:        userClient,
		orchestraProducer: orchestraProducer,
		provider:          provider,
	}
}

func (u *Usecase) CreateUserMessaging(ctx context.Context, ge event.GlobalEvent[dto.UserCreateRequest, any]) error {

	response, err := u.provider.CreateUser(ctx, &ge.Payload.Request)

	var gevent event.GlobalEvent[dto.UserCreateRequest, any]

	basePayload := event.BasePayload[dto.UserCreateRequest, any]{
		Request: ge.Payload.Request,
	}

	if err != nil {
		basePayload.Response = err

		gevent = event.NewGlobalEvent[dto.UserCreateRequest, any](
			"create",
			"error",
			"user_creation_failed",
			basePayload,
		)
		if response.Error != "" {
			gevent.StatusCode = response.StatusCode
		} else {
			gevent.StatusCode = 500
		}
	} else {
		basePayload.Response = response.Data
		gevent = event.NewGlobalEvent[dto.UserCreateRequest, any](
			"create",
			"success",
			"user_created",
			basePayload,
		)

		gevent.StatusCode = response.StatusCode
	}

	gevent.EventID = ge.EventID
	gevent.InstanceID = ge.InstanceID
	gevent.EventType = ge.EventType

	log.Println("check state:", gevent.State)

	bytes, jsonErr := gevent.ToJSON()
	if jsonErr != nil {
		return fmt.Errorf("failed to convert event to JSON: %w", jsonErr)
	}

	sendErr := u.orchestraProducer.SendMessage(uuid.New().String(), bytes)
	if sendErr != nil {
		return fmt.Errorf("failed to send message: %w", sendErr)
	}

	if err != nil {
		return fmt.Errorf("user creation failed: %s", err)
	}

	return nil
}

func (u *Usecase) UpdateUserMessaging(ctx context.Context, ge event.GlobalEvent[dto.UpdateBankIDRequest, any]) error {

	response, err := u.provider.UpdateUser(ctx, &dto.UpdateBankIDRequest{
		Username:      ge.Payload.Request.Username,
		AccountBankID: ge.Payload.Request.AccountBankID,
	})

	var gevent event.GlobalEvent[dto.UpdateBankIDRequest, any]

	basePayload := event.BasePayload[dto.UpdateBankIDRequest, any]{
		Request: ge.Payload.Request,
	}

	if err != nil {
		basePayload.Response = err

		gevent = event.NewGlobalEvent[dto.UpdateBankIDRequest, any](
			"update",
			"error",
			"user_update_failed",
			basePayload,
		)
		if response.Error != "" {
			gevent.StatusCode = response.StatusCode
		} else {
			gevent.StatusCode = 500
		}
	} else {
		basePayload.Response = response.Data
		gevent = event.NewGlobalEvent[dto.UpdateBankIDRequest, any](
			"update",
			"success",
			"user_bankid_updated",
			basePayload,
		)

		gevent.StatusCode = response.StatusCode
	}

	gevent.EventID = ge.EventID
	gevent.InstanceID = ge.InstanceID
	gevent.EventType = ge.EventType

	log.Println("check state:", gevent.State)

	bytes, jsonErr := gevent.ToJSON()
	if jsonErr != nil {
		return fmt.Errorf("failed to convert event to JSON: %w", jsonErr)
	}

	sendErr := u.orchestraProducer.SendMessage(uuid.New().String(), bytes)

	if sendErr != nil {
		return fmt.Errorf("failed to send message: %w", sendErr)
	}

	if err != nil {
		return fmt.Errorf("user update failed: %s", err)
	}

	return nil
}

func (u *Usecase) UserDetailMessaging(ctx context.Context, ge event.GlobalEvent[dto.UserValidateRequest, any]) error {

	response, err := u.provider.GetUserDetail(ctx, &dto.UserValidateRequest{
		Username: ge.Payload.Request.Username,
	})

	var gevent event.GlobalEvent[dto.UserValidateRequest, any]
	basePayload := event.BasePayload[dto.UserValidateRequest, any]{
		Request: ge.Payload.Request,
	}

	if err != nil {
		basePayload.Response = err

		gevent = event.NewGlobalEvent[dto.UserValidateRequest, any](
			"get",
			"error",
			"user_validation_failed",
			basePayload,
		)

		if response.Error != "" {
			gevent.StatusCode = response.StatusCode
		} else {
			gevent.StatusCode = 500
		}

	} else {
		basePayload.Response = response.Data
		gevent = event.NewGlobalEvent[dto.UserValidateRequest, any](
			"get",
			"success",
			"user_validation_success",
			basePayload,
		)
		gevent.StatusCode = response.StatusCode
	}

	gevent.EventID = ge.EventID
	gevent.InstanceID = ge.InstanceID
	gevent.EventType = ge.EventType

	log.Println("check state:", gevent.State)

	bytes, jsonErr := gevent.ToJSON()
	if jsonErr != nil {
		return fmt.Errorf("failed to convert event to JSON: %w", jsonErr)
	}

	sendErr := u.orchestraProducer.SendMessage(uuid.New().String(), bytes)
	if sendErr != nil {
		return fmt.Errorf("failed to send message: %w", sendErr)
	}

	if err != nil {
		return fmt.Errorf("user validation failed: %s", err)
	}

	return nil
}
