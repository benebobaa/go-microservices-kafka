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
	return &OrderUsecase{queries: queries, orchestraProducer: producer}
}

func (oc *OrderUsecase) CreateOrder(ctx context.Context, order *dto.OrderRequest) (*sqlc.Order, error) {

	orderCreated, err := oc.queries.CreateOrder(ctx, sqlc.CreateOrderParams{
		OrderRefID: fmt.Sprintf("%s-%s", "TOKPED", uuid.New()),
		CustomerID: order.CustomerID,
		Username:   order.Username,
		ProductID:  order.ProductID,
		Quantity:   order.Quantity,
		Status:     dto.PROCESSING.String(),
	})

	if err != nil {
		return nil, err
	}

	orderEvent := event.NewGlobalEvent(
		"create",
		"success",
		"order_process",
		orderCreated,
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

func (oc *OrderUsecase) CancelOrder(ctx context.Context, req *dto.OrderCancelRequest) (*sqlc.Order, error) {
	order, err := oc.queries.FindOrderByID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	if req.Username != order.Username {
		return nil, ErrUnauthorizeCancelOrder
	}

	if order.Status != dto.COMPLETE.String() {
		return nil, ErrCannotCancelOrder
	}

	cancelledOrder, err := oc.queries.UpdateOrder(ctx, sqlc.UpdateOrderParams{
		Status:      dto.CANCEL_PROCESSING.String(),
		TotalAmount: order.TotalAmount,
		OrderRefID:  order.OrderRefID,
	})

	if err != nil {
		return nil, err
	}

	orderEvent := event.NewGlobalEvent(
		"update",
		"success",
		"order_cancel",
		cancelledOrder,
	)

	orderEvent.State = "ORDER_CANCELLED"

	bytes, err := orderEvent.ToJSON()
	if err != nil {
		return nil, err
	}

	err = oc.orchestraProducer.SendMessage(uuid.New().String(), bytes)
	if err != nil {
		return nil, err
	}

	return &cancelledOrder, nil
}

func (oc *OrderUsecase) UpdateOrder(ctx context.Context, updateOrder *dto.OrderUpdateRequest) (*sqlc.Order, error) {

	count, err := oc.queries.CountByID(ctx, updateOrder.OrderRefID)

	if err != nil {
		return nil, err
	}

	if count < 1 {
		return nil, errors.New("order id not found")
	}

	updatedOrder, err := oc.queries.UpdateOrder(ctx, sqlc.UpdateOrderParams{
		Status:      updateOrder.Status,
		TotalAmount: updateOrder.TotalAmount,
		OrderRefID:  updateOrder.OrderRefID,
	})

	if err != nil {
		return nil, err
	}

	orderEvent := event.NewGlobalEvent(
		"update",
		"success",
		"order_updated",
		updatedOrder,
	)

	bytes, err := orderEvent.ToJSON()

	if err != nil {
		return nil, err
	}

	err = oc.orchestraProducer.SendMessage(uuid.New().String(), bytes)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
