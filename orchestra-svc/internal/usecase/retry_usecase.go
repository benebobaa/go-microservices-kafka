package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"orchestra-svc/internal/dto"
	"orchestra-svc/internal/dto/event"
	"orchestra-svc/internal/repository/sqlc"
	"orchestra-svc/pkg/producer"
	"time"
)

type RetryUsecase struct {
	queries  sqlc.Store
	producer *producer.KafkaProducer
	oc       *OrchestraUsecase
}

func NewRetryUsecase(
	queries sqlc.Store,
	producer *producer.KafkaProducer,
	oc *OrchestraUsecase,
) *RetryUsecase {
	return &RetryUsecase{
		queries:  queries,
		producer: producer,
		oc:       oc,
	}
}

func (r *RetryUsecase) ProductQuantityRetry(ctx context.Context, req *dto.ProductQuantityRetryRequest) (*event.GlobalEvent[dto.ProductReserveRequest, any], error) {

	insStep, err := r.queries.FindWorkflowInstanceStepsByEventIDAndInsID(ctx, sqlc.FindWorkflowInstanceStepsByEventIDAndInsIDParams{
		EventID:            req.EventID,
		WorkflowInstanceID: req.InstanceID,
	})

	if err != nil {
		return nil, err
	}

	if insStep.Status != dto.ERROR.String() {
		return nil, errors.New("step is not in failed state")
	}

	if !insStep.EventMessage.Valid {
		return nil, errors.New("event message is not valid")
	}

	eventMsg, err := event.FromJSON[dto.ProductReserveRequest, any]([]byte(insStep.EventMessage.String))

	if err != nil {
		return nil, err
	}

	gevent := event.NewGlobalEvent[dto.ProductReserveRequest, any](
		"retry", "success", eventMsg.Payload)

	gevent.State = "product_retry"
	gevent.StatusCode = 200
	gevent.Payload.Request.Quantity = req.Quantity

	gevent.EventType = eventMsg.EventType
	gevent.InstanceID = eventMsg.InstanceID
	gevent.EventID = eventMsg.EventID

	bytes, err := gevent.ToJSON()
	if err != nil {
		return nil, err
	}

	err = r.queries.UpdateWorkflowInstanceStep(ctx, sqlc.UpdateWorkflowInstanceStepParams{
		Status:       dto.IN_PROGRESS.String(),
		EventMessage: sql.NullString{String: string(bytes), Valid: true},
		StartedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		EventID:      gevent.EventID,
	})

	if err != nil {
		return nil, err
	}

	err = r.producer.SendMessage(insStep.Topic, uuid.New().String(), bytes)

	if err != nil {
		return nil, err
	}

	return &gevent, nil
}

func (r *RetryUsecase) RetryFailedInstanceStep(ctx context.Context, req *dto.RetryRequest) (*event.GlobalEvent[any, any], error) {

	insStep, err := r.queries.FindWorkflowInstanceStepsByEventIDAndInsID(ctx, sqlc.FindWorkflowInstanceStepsByEventIDAndInsIDParams{
		EventID:            req.EventID,
		WorkflowInstanceID: req.InstanceID,
	})

	if err != nil {
		return nil, err
	}

	if insStep.Status != dto.ERROR.String() {
		return nil, errors.New("step is not in failed state")
	}

	if !insStep.EventMessage.Valid {
		return nil, errors.New("event message is not valid")
	}

	eventMsg, err := event.FromJSON[any, any]([]byte(insStep.EventMessage.String))

	if err != nil {
		return nil, err
	}

	gevent := event.NewGlobalEvent[any, any](
		"retry", "success", eventMsg.Payload)

	gevent.Payload.Request = eventMsg.Payload.Request
	gevent.State = eventMsg.State
	gevent.EventType = eventMsg.EventType
	gevent.InstanceID = eventMsg.InstanceID
	gevent.EventID = eventMsg.EventID
	gevent.StatusCode = 200

	bytes, err := gevent.ToJSON()
	if err != nil {
		return nil, err
	}

	err = r.queries.UpdateWorkflowInstanceStep(ctx, sqlc.UpdateWorkflowInstanceStepParams{
		Status:       dto.IN_PROGRESS.String(),
		EventMessage: sql.NullString{String: string(bytes), Valid: true},
		StartedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		EventID:      gevent.EventID,
	})

	if err != nil {
		return nil, err
	}

	err = r.producer.SendMessage(insStep.Topic, uuid.New().String(), bytes)

	if err != nil {
		return nil, err
	}

	return &gevent, nil
}
