package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"product-svc/internal/dto"
	"product-svc/internal/dto/event"
	"product-svc/internal/provider"
	"product-svc/internal/usecase"
	"product-svc/pkg/producer"
)

func TestUsecase_ReserveProductMessaging(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductProvider := provider.NewMockProductProvider(ctrl)
	mockKafkaProducer := producer.NewMockKafkaProducer(ctrl)
	u := usecase.NewUsecase(mockProductProvider, mockKafkaProducer)

	ctx := context.Background()
	ge := event.GlobalEvent[dto.ProductRequest, any]{
		EventID:    "event-id",
		InstanceID: "instance-id",
		EventType:  "event-type",
		Payload: event.BasePayload[dto.ProductRequest, any]{
			Request: dto.ProductRequest{
				ProductID: "product-id",
				Quantity:  1,
			},
		},
	}

	t.Run("successful reservation", func(t *testing.T) {
		mockProductProvider.EXPECT().ReserveProduct(ctx, &dto.ProductRequest{
			ProductID: "product-id",
			Quantity:  1,
		}).Return(&dto.BaseResponse[dto.ProductResponse]{Data: &dto.ProductResponse{Id: "product-id", Name: "product-name", Quantity: 1, Price: 10.0, Amount: 10.0}, StatusCode: 200}, nil)

		mockKafkaProducer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil)

		err := u.ReserveProductMessaging(ctx, ge)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("message sending failure", func(t *testing.T) {
		mockProductProvider.EXPECT().ReserveProduct(ctx, &dto.ProductRequest{
			ProductID: "product-id",
			Quantity:  1,
		}).Return(&dto.BaseResponse[dto.ProductResponse]{Data: &dto.ProductResponse{Id: "product-id", Name: "product-name", Quantity: 1, Price: 10.0, Amount: 10.0}, StatusCode: 200}, nil)

		mockKafkaProducer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(errors.New("send error"))

		err := u.ReserveProductMessaging(ctx, ge)
		if err == nil || err.Error() != "failed to send message: send error" {
			t.Errorf("expected send error, got %v", err)
		}
	})
}

func TestUsecase_ReleaseProductMessaging(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductProvider := provider.NewMockProductProvider(ctrl)
	mockKafkaProducer := producer.NewMockKafkaProducer(ctrl)
	u := usecase.NewUsecase(mockProductProvider, mockKafkaProducer)

	ctx := context.Background()
	ge := event.GlobalEvent[dto.ProductRequest, any]{
		EventID:    "event-id",
		InstanceID: "instance-id",
		EventType:  "event-type",
		Payload: event.BasePayload[dto.ProductRequest, any]{
			Request: dto.ProductRequest{
				ProductID: "product-id",
				Quantity:  1,
			},
		},
	}

	t.Run("successful release", func(t *testing.T) {
		mockProductProvider.EXPECT().ReleaseProduct(ctx, &dto.ProductRequest{
			ProductID: "product-id",
			Quantity:  1,
		}).Return(&dto.BaseResponse[dto.ProductResponse]{Data: &dto.ProductResponse{Id: "product-id", Name: "product-name", Quantity: 1, Price: 10.0, Amount: 10.0}, StatusCode: 200}, nil)

		mockKafkaProducer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil)

		err := u.ReleaseProductMessaging(ctx, ge)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	//t.Run("release response error", func(t *testing.T) {
	//	mockProductProvider.EXPECT().ReleaseProduct(ctx, &dto.ProductRequest{
	//		ProductID: "product-id",
	//		Quantity:  1,
	//	}).Return(&dto.BaseResponse[dto.ProductResponse]{Error: "response error", StatusCode: 400}, nil)
	//
	//	mockKafkaProducer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(nil)
	//
	//	err := u.ReleaseProductMessaging(ctx, ge)
	//	if err == nil || err.Error() != "product release failed: response error" {
	//		t.Errorf("expected response error, got %v", err)
	//	}
	//})

	t.Run("message sending failure", func(t *testing.T) {
		mockProductProvider.EXPECT().ReleaseProduct(ctx, &dto.ProductRequest{
			ProductID: "product-id",
			Quantity:  1,
		}).Return(&dto.BaseResponse[dto.ProductResponse]{Data: &dto.ProductResponse{Id: "product-id", Name: "product-name", Quantity: 1, Price: 10.0, Amount: 10.0}, StatusCode: 200}, nil)

		mockKafkaProducer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return(errors.New("send error"))

		err := u.ReleaseProductMessaging(ctx, ge)
		if err == nil || err.Error() != "failed to send message: send error" {
			t.Errorf("expected send error, got %v", err)
		}
	})
}
