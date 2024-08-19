package usecase

import (
	"context"
	"fmt"
	retryit "github.com/benebobaa/retry-it"
	"github.com/google/uuid"
	"log"
	"time"
	"user-svc/internal/dto"
	"user-svc/internal/dto/event"
	"user-svc/internal/interfaces"
)

type Usecase struct {
	userClient        interfaces.Client
	orchestraProducer interfaces.Producer
}

func NewUsecase(userClient interfaces.Client, orchestraProducer interfaces.Producer) *Usecase {
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

func (u *Usecase) ValidateUser(ctx context.Context, request *dto.UserValidateRequest) (*dto.BaseResponse[dto.UserResponse], *dto.ErrorResponse) {

	var response dto.BaseResponse[dto.UserResponse]
	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		log.Println("retrying: ", counter)
		return u.userClient.GET(
			ctx,
			fmt.Sprintf("/users/%s", request.Username),
			request,
			&response,
		)
	}, retryit.WithInitialDelay(500*time.Millisecond))

	log.Println("response:", response)

	if err != nil {
		log.Println("error:", err.Error())
		return &response, &dto.ErrorResponse{
			Message: err.Error(),
		}
	}

	if response.StatusCode != 200 {
		return &response, &dto.ErrorResponse{
			Message: response.Error,
		}
	}

	return &response, nil
}

//func
