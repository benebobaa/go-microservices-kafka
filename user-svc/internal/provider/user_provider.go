package provider

import (
	"context"
	"fmt"
	retryit "github.com/benebobaa/retry-it"
	"log"
	"time"
	"user-svc/internal/dto"
	"user-svc/internal/interfaces"
)

type UserProvider interface {
	GetUserDetail(ctx context.Context, request *dto.UserValidateRequest) (*dto.BaseResponse[dto.UserResponse], *dto.ErrorResponse)
	UpdateUser(ctx context.Context, request *dto.UpdateBankIDRequest) (*dto.BaseResponse[dto.UserResponse], *dto.ErrorResponse)
	CreateUser(ctx context.Context, request *dto.UserCreateRequest) (*dto.BaseResponse[dto.UserResponse], *dto.ErrorResponse)
}

type UserProviderImpl struct {
	client interfaces.Client
}

func NewUserProvider(client interfaces.Client) UserProvider {
	return &UserProviderImpl{
		client: client,
	}
}

func (u *UserProviderImpl) GetUserDetail(ctx context.Context, request *dto.UserValidateRequest) (*dto.BaseResponse[dto.UserResponse], *dto.ErrorResponse) {

	var response dto.BaseResponse[dto.UserResponse]
	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		log.Println("retrying: ", counter)
		return u.client.GET(
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

func (u *UserProviderImpl) UpdateUser(ctx context.Context, request *dto.UpdateBankIDRequest) (*dto.BaseResponse[dto.UserResponse], *dto.ErrorResponse) {

	var response dto.BaseResponse[dto.UserResponse]
	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		log.Println("retrying: ", counter)
		return u.client.PATCH(
			ctx,
			fmt.Sprintf("/users/%s", request.Username),
			request,
			&response,
		)
	}, retryit.WithInitialDelay(500*time.Millisecond))

	if err != nil {
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

func (u *UserProviderImpl) CreateUser(ctx context.Context, request *dto.UserCreateRequest) (*dto.BaseResponse[dto.UserResponse], *dto.ErrorResponse) {

	var response dto.BaseResponse[dto.UserResponse]
	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		log.Println("retrying: ", counter)
		return u.client.POST(
			ctx,
			"/users",
			request,
			&response,
		)
	}, retryit.WithInitialDelay(500*time.Millisecond))

	if err != nil {
		return &response, &dto.ErrorResponse{
			Message: err.Error(),
		}
	}

	if response.StatusCode != 201 {
		return &response, &dto.ErrorResponse{
			Message: response.Error,
		}
	}

	return &response, nil
}
