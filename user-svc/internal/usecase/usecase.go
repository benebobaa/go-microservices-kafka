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

func (u *Usecase) ValidateUserMessaging(ctx context.Context, ge event.GlobalEvent[dto.UserValidateRequest]) error {

	response, err := u.ValidateUser(ctx, &dto.UserValidateRequest{
		Username: ge.Payload.Username,
	})

	if err != nil {
		return err
	}

	log.Println("response data: ", response.Data)

	gevent := event.NewGlobalEvent(
		"get",
		"success",
		"user_validated",
		"user_validation_success",
		response.Data,
	)

	gevent.EventID = ge.EventID
	gevent.InstanceID = ge.InstanceID
	gevent.EventType = ge.EventType
	gevent.StatusCode = response.StatusCode
	bytes, err := gevent.ToJSON()

	log.Println("gevent: ", gevent)
	if err != nil {
		return err
	}
	err = u.orchestraProducer.SendMessage(uuid.New().String(), bytes)

	if err != nil {
		return err
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
		return nil, err
	}

	log.Println("response userClient: ", response)

	return &response, nil
}
