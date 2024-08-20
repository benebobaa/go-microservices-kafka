package interfaces

import (
	"context"
	"product-svc/internal/dto"
	"product-svc/internal/dto/event"
)

type ProductUsecase interface {
	ReserveProductMessaging(ctx context.Context, ge event.GlobalEvent[dto.ProductRequest, any]) error
	ReleaseProductMessaging(ctx context.Context, ge event.GlobalEvent[dto.ProductRequest, any]) error
}
