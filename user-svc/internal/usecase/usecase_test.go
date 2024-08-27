package usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"testing"
	"user-svc/internal/dto"
	"user-svc/internal/dto/event"
	"user-svc/internal/provider"
	"user-svc/pkg/producer"
)

var orchProducer *producer.KafkaProducer

func init() {
	orchProducer, _ = producer.NewKafkaProducer([]string{"localhost:29092"}, "orchestra-topic-test")
}

func TestUsecase_CreateUserMessaging(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := provider.NewMockUserProvider(ctrl)
	userUsecase := NewUsecase(orchProducer, mockProvider)
	ge := event.GlobalEvent[dto.UserCreateRequest, any]{
		EventID:    "event-id",
		InstanceID: "instance-id",
		EventType:  "event-type",
		Payload: event.BasePayload[dto.UserCreateRequest, any]{
			Request: dto.UserCreateRequest{
				Username: "beneboba",
				Email:    "beneboba@gmail.com",
			},
		},
	}
	t.Run("successful user creation", func(t *testing.T) {

		mockProvider.EXPECT().
		//mockProvider.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, nil)

		err := userUsecase.CreateUserMessaging(context.Background(), ge)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}
