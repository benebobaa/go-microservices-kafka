package usecase

import (
	"context"
	"errors"
	"fmt"
	"user-svc/internal/dto"
	"user-svc/internal/dto/event"
	"user-svc/pkg/http_client"
	"user-svc/pkg/producer"

	"github.com/google/uuid"
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

func (u *Usecase) ValidateUser(ctx context.Context, request *dto.UserValidateRequest) (*dto.UserResponse, error) {

	var response dto.BaseResponse[dto.UserResponse]

	gEvent := event.NewGlobalEvent(
		"get",
		"success",
		"user_valid",
		response,
	)

	defer func() {
		bytes, _ := gEvent.ToJSON()
		u.orchestraProducer.SendMessage(uuid.New().String(), bytes)
	}()

	err := u.userClient.GET(
		fmt.Sprintf("/users/%s", request.Username),
		nil,
		&response,
	)

	if err != nil || response.Error != "" {
		gEvent.Status = "failed"
		gEvent.EventType = "user_not_valid"
		return nil, errors.Join(err, errors.New(response.Error))
	}

	bytes, _ := gEvent.ToJSON()
	err = u.orchestraProducer.SendMessage(uuid.New().String(), bytes)

	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}
