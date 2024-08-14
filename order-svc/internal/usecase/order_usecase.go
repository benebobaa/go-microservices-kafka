package usecase

import (
	"context"
	"errors"
	"fmt"
	"order-svc/internal/dto"
	"order-svc/internal/dto/event"
	"order-svc/internal/repository/sqlc"
	"order-svc/pkg/producer"

	"github.com/google/uuid"
)

var ErrUnauthorizeCancelOrder = errors.New("user unauthorize to cancel order")
var ErrCannotCancelOrder = errors.New("cannot cancel uncomplete order")

type OrderUsecase struct {
	queries           sqlc.Querier
	orchestraProducer *producer.KafkaProducer
}

func NewOrderUsecase(queries sqlc.Querier, producer *producer.KafkaProducer) *OrderUsecase {
	return &OrderUsecase{
		queries:           queries,
		orchestraProducer: producer,
	}
}

func (oc *OrderUsecase) CreateOrder(ctx context.Context, req *dto.OrderRequest) (*sqlc.Order, error) {

	orderCreated, err := oc.queries.CreateOrder(ctx, sqlc.CreateOrderParams{
		RefID:      fmt.Sprintf("%s-%s", "TOKPED", uuid.New()),
		CustomerID: req.CustomerID,
		Username:   req.Username,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		Status:     dto.PROCESSING.String(),
	})

	if err != nil {
		return nil, err
	}

	basePayload := event.BasePayload[dto.OrderRequest, sqlc.Order]{
		Request:  *req,
		Response: orderCreated,
	}

	orderEvent := event.NewGlobalEvent(
		"create",
		"success",
		"order_process",
		basePayload,
	)
	orderEvent.State = "order_created"
	orderEvent.StatusCode = 201
	bytes, err := orderEvent.ToJSON()

	if err != nil {
		return nil, err
	}

	err = oc.orchestraProducer.SendMessage(uuid.New().String(), bytes)

	if err != nil {
		return nil, err
	}

	return &orderCreated, nil
}

func (oc *OrderUsecase) CancelOrder(ctx context.Context, username string, id int) (*sqlc.Order, error) {
	order, err := oc.queries.FindOrderByID(ctx, int32(id))

	if err != nil {
		return nil, err
	}

	if order.Status != dto.COMPLETE.String() {
		return nil, ErrCannotCancelOrder
	}

	if username != order.Username {
		return nil, ErrUnauthorizeCancelOrder
	}

	return oc.UpdateOrderMessaging(ctx, "order_cancel_process", &dto.OrderUpdateRequest{
		RefID:  order.RefID,
		Status: dto.CANCEL_PROCESSING.String(),
		Amount: order.TotalAmount,
	})
}

func (oc *OrderUsecase) UpdateOrderMessaging(ctx context.Context, eventType string, req *dto.OrderUpdateRequest) (*sqlc.Order, error) {

	updatedOrder, err := oc.queries.UpdateOrder(ctx, sqlc.UpdateOrderParams{
		Status:      req.Status,
		TotalAmount: req.Amount,
		RefID:       req.RefID,
	})

	if err != nil {
		return nil, err
	}

	basePayload := event.BasePayload[dto.OrderUpdateRequest, sqlc.Order]{
		Request:  *req,
		Response: updatedOrder,
	}

	orderEvent := event.NewGlobalEvent(
		"create",
		"success",
		"order_process",
		basePayload,
	)

	orderEvent.State = "order_created"
	orderEvent.StatusCode = 200

	if eventType != "" {
		orderEvent.EventType = eventType
		orderEvent.State = "order_cancelled"
	}
	bytes, err := orderEvent.ToJSON()

	if err != nil {
		return nil, err
	}

	err = oc.orchestraProducer.SendMessage(uuid.New().String(), bytes)

	if err != nil {
		return nil, err
	}

	return &updatedOrder, nil
}
