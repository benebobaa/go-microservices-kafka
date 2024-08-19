package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
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
		"order_created",
		event.ORDER_PROCESS.String(),
		basePayload,
	)
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
	order, err := oc.queries.FindOrderByID(ctx, int32(req.OrderID))
	if err != nil {
		return nil, err
	}

	if order.Status != dto.COMPLETE.String() {
		return nil, ErrCannotCancelOrder
	}

	if req.Username != order.Username {
		return nil, ErrUnauthorizeCancelOrder
	}

	reqUpdate := &dto.OrderUpdateRequest{
		RefID:     order.RefID,
		Amount:    order.Amount.Float64,
		Quantity:  order.Quantity,
		Status:    dto.CANCEL_PROCESSING.String(),
		EventType: event.ORDER_CANCEL_PROCESS.String(),
	}

	updatedOrder, err := oc.UpdateOrder(ctx, reqUpdate)
	if err != nil {
		return nil, err
	}

	basePayload := event.BasePayload[dto.OrderUpdateRequest, sqlc.Order]{
		Request:  *reqUpdate,
		Response: *updatedOrder,
	}

	orderEvent := event.NewGlobalEvent(
		"update",
		"success",
		"order_cancel",
		event.ORDER_CANCEL_PROCESS.String(),
		basePayload,
	)

	orderEvent.StatusCode = 200

	bytes, err := orderEvent.ToJSON()
	if err != nil {
		return nil, err
	}

	err = oc.orchestraProducer.SendMessage(uuid.New().String(), bytes)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (oc *OrderUsecase) UpdateOrderMessaging(ctx context.Context, req event.GlobalEvent[dto.OrderUpdateRequest, any]) error {

	order, err := oc.queries.FindOrderByRefID(ctx, req.Payload.Request.RefID)

	if err != nil {
		log.Println("error find order by refID: ", err)
		return err
	}

	var updateReq dto.OrderUpdateRequest

	if req.EventType == event.ORDER_PROCESS.String() {
		updateReq = dto.OrderUpdateRequest{
			RefID:     order.RefID,
			Amount:    req.Payload.Request.Amount,
			Status:    req.Payload.Request.Status,
			Quantity:  req.Payload.Request.Quantity,
			EventType: req.EventType,
		}
	} else {
		updateReq = dto.OrderUpdateRequest{
			RefID:     order.RefID,
			Amount:    order.Amount.Float64,
			Status:    req.Payload.Request.Status,
			Quantity:  order.Quantity,
			EventType: req.EventType,
		}
	}

	updatedOrder, err := oc.UpdateOrder(ctx, &updateReq)

	if err != nil {
		return err
	}

	basePayload := event.BasePayload[dto.OrderUpdateRequest, sqlc.Order]{
		Request:  req.Payload.Request,
		Response: *updatedOrder,
	}

	orderEvent := event.NewGlobalEvent(
		"update",
		"success",
		"order_updated",
		req.EventType,
		basePayload,
	)

	orderEvent.EventID = req.EventID
	orderEvent.InstanceID = req.InstanceID
	orderEvent.StatusCode = 200

	bytes, err := orderEvent.ToJSON()
	if err != nil {
		return err
	}

	err = oc.orchestraProducer.SendMessage(uuid.New().String(), bytes)
	if err != nil {
		return err
	}

	return nil
}

func (oc *OrderUsecase) UpdateOrder(ctx context.Context, req *dto.OrderUpdateRequest) (*sqlc.Order, error) {

	log.Println("req updateOrder: ", req)

	updatedOrder, err := oc.queries.UpdateOrder(ctx, sqlc.UpdateOrderParams{
		Status: req.Status,
		Amount: sql.NullFloat64{
			Float64: req.Amount,
			Valid:   true,
		},
		Quantity: req.Quantity,
		RefID:    req.RefID,
	})

	if err != nil {
		return nil, err
	}

	return &updatedOrder, nil
}

func (oc *OrderUsecase) FindAllOrder(ctx context.Context, username string) ([]sqlc.Order, error) {
	orders, err := oc.queries.FindOrdersByUsername(ctx, username)

	if err != nil {
		return nil, err
	}

	return orders, nil
}
