package provider

import (
	"context"
	retryit "github.com/benebobaa/retry-it"
	"log"
	"product-svc/internal/dto"
	"product-svc/internal/interfaces"
	"time"
)

type ProductProvider interface {
	ReserveProduct(ctx context.Context, req *dto.ProductRequest) (*dto.BaseResponse[dto.ProductResponse], *dto.ErrorResponse)
	ReleaseProduct(ctx context.Context, req *dto.ProductRequest) (*dto.BaseResponse[dto.ProductResponse], *dto.ErrorResponse)
}

type ProductProviderImpl struct {
	client interfaces.HtppRequest
}

func NewProductProviderImpl(client interfaces.HtppRequest) ProductProvider {
	return &ProductProviderImpl{client: client}
}

func (u *ProductProviderImpl) ReserveProduct(ctx context.Context, req *dto.ProductRequest) (*dto.BaseResponse[dto.ProductResponse], *dto.ErrorResponse) {
	var response dto.BaseResponse[dto.ProductResponse]

	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		log.Println("retrying reserve: ", counter)
		return u.client.POST(ctx, "/reserve", req, &response)
	}, retryit.WithInitialDelay(500*time.Millisecond))

	//err := u.client.POST(ctx, "/reserve", req, &response)

	if err != nil {
		return &response, &dto.ErrorResponse{Error: err.Error()}
	}

	if response.StatusCode != 200 {
		return &response, &dto.ErrorResponse{Error: response.Error}
	}

	return &response, nil
}

func (u *ProductProviderImpl) ReleaseProduct(ctx context.Context, req *dto.ProductRequest) (*dto.BaseResponse[dto.ProductResponse], *dto.ErrorResponse) {
	var response dto.BaseResponse[dto.ProductResponse]

	counter := 0
	err := retryit.Do(ctx, func(ctx context.Context) error {
		counter++
		log.Println("retrying release: ", counter)
		return u.client.POST(ctx, "/release", req, &response)
	}, retryit.WithInitialDelay(500*time.Millisecond))
	//err := u.client.POST(ctx, "/release", req, &response)

	if err != nil {
		return &response, &dto.ErrorResponse{Error: err.Error()}
	}

	if response.StatusCode != 200 {
		return &response, &dto.ErrorResponse{Error: response.Error}
	}

	return &response, nil
}
