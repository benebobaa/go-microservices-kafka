package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"user-svc/internal/dto"
	"user-svc/internal/dto/event"
	"user-svc/pkg/http_client"
	"user-svc/pkg/producer"
)

type Usecase struct {
	userClient        *http_client.UserClient
	orchestraProducer *producer.KafkaProducer
}

func NewUsecase(userClient *http_client.UserClient, orchestraProducer *producer.KafkaProducer) *Usecase {
	return &Usecase{
		userClient:        userClient,
		orchestraProducer: orchestraProducer,
	}
}

func (u *Usecase) ValidateUserMessaging(ctx context.Context, ge event.GlobalEvent[dto.UserValidateRequest, any]) error {

	response, err := u.ValidateUser(ctx, &dto.UserValidateRequest{
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
	} else {
		basePayload.Response = response.Data
		gevent = event.NewGlobalEvent[dto.UserValidateRequest, any](
			"get",
			"success",
			"user_validation_success",
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
		return fmt.Errorf("user validation failed: %w", err)
	}

	return nil
}

func (u *Usecase) ValidateUser(ctx context.Context, request *dto.UserValidateRequest) (*dto.BaseResponse[dto.UserResponse], error) {

	var response dto.BaseResponse[dto.UserResponse]

	err := u.userClient.GET(
		fmt.Sprintf("/users/%s", request.Username),
		request,
		&response,
	)

	if err != nil {
		return &response, err
	}

	return &response, nil
}
