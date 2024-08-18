package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"user-svc/internal/dto"
	"user-svc/internal/dto/event"
)

// Mocking UserClient
type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) GET(ctx context.Context, url string, request interface{}, response interface{}) error {
	args := m.Called(ctx, url, request, response)
	return args.Error(0)
}

// Mocking KafkaProducer
type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) SendMessage(key string, message []byte) error {
	args := m.Called(key, message)
	return args.Error(0)
}

func TestValidateUserMessaging(t *testing.T) {
	userClient := new(MockUserClient)
	kafkaProducer := new(MockKafkaProducer)
	usecase := NewUsecase(userClient, kafkaProducer)

	// Test setup
	ge := event.GlobalEvent[dto.UserValidateRequest, any]{
		EventID:    "event-id",
		InstanceID: "instance-id",
		EventType:  "event-type",
		Payload: event.BasePayload[dto.UserValidateRequest, any]{
			Request: dto.UserValidateRequest{
				Username: "testuser",
			},
		},
	}

	// Successful case
	userClient.On("GET", mock.Anything, "/users/testuser", mock.Anything).Run(func(args mock.Arguments) {
		response := args.Get(3).(*dto.BaseResponse[dto.UserResponse])
		response.StatusCode = 200
	}).Return(nil)

	kafkaProducer.On("SendMessage", mock.Anything, mock.Anything).Return(nil)

	err := usecase.ValidateUserMessaging(context.Background(), ge)
	assert.NoError(t, err)

	// Test error case
	userClient.On("GET", mock.Anything, "/users/testuser", mock.Anything).Return(errors.New("network error"))

	kafkaProducer.On("SendMessage", mock.Anything, mock.Anything).Return(nil)

	err = usecase.ValidateUserMessaging(context.Background(), ge)
	assert.Error(t, err)

	// Verify method calls
	userClient.AssertExpectations(t)
	kafkaProducer.AssertExpectations(t)
}

func TestValidateUser(t *testing.T) {
	userClient := new(MockUserClient)
	usecase := NewUsecase(userClient, nil)

	// Successful case
	userClient.On("GET", mock.Anything, "/users/testuser", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		response := args.Get(3).(*dto.BaseResponse[dto.UserResponse])
		response.StatusCode = 200
	}).Return(nil)

	response, err := usecase.ValidateUser(context.Background(), &dto.UserValidateRequest{Username: "testuser"})
	assert.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	// Error case
	userClient.On("GET", mock.Anything, "/users/testuser", mock.Anything, mock.Anything).Return(errors.New("network error"))

	response, err = usecase.ValidateUser(context.Background(), &dto.UserValidateRequest{Username: "testuser"})
	assert.Error(t, err)
	assert.Empty(t, response.StatusCode)
}
