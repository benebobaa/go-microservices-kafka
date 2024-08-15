package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"orchestra-svc/internal/dto"
	"orchestra-svc/internal/dto/event"
	"orchestra-svc/internal/repository/cache"
	"orchestra-svc/internal/repository/sqlc"
	"orchestra-svc/pkg"
	"orchestra-svc/pkg/producer"
	"time"
)

type OrchestraUsecase struct {
	queries  sqlc.Querier
	producer *producer.KafkaProducer
	cache    *cache.PayloadCacher
}

func NewOrchestraUsecase(q sqlc.Querier, p *producer.KafkaProducer, c *cache.PayloadCacher) *OrchestraUsecase {
	return &OrchestraUsecase{
		queries:  q,
		producer: p,
		cache:    c,
	}
}

func (o *OrchestraUsecase) ProcessWorkflow(ctx context.Context, eventMsg event.GlobalEvent[any, any]) error {

	cachePayload, err := o.getCachePayload(eventMsg.InstanceID, eventMsg.Source, eventMsg.Payload.Response)
	if err != nil {
		return err
	}

	wf, err := o.queries.FindWorkflowByType(ctx, eventMsg.EventType)
	if err != nil {
		return fmt.Errorf("find workflow: %w", err)
	}

	err = o.handleInstanceStep(ctx, eventMsg)
	if err != nil {
		log.Println("Error handling instance step: ", err)
		return err
	}

	instance, err := o.getOrCreateWorkflowInstance(ctx, eventMsg, wf)
	if err != nil {
		log.Println("Error get or create workflow instance: ", err)
		return err
	}

	return o.processSteps(ctx, eventMsg, instance, cachePayload)
}

func (o *OrchestraUsecase) getCachePayload(instanceID string, source string, response any) (map[string]any, error) {
	cachePayload, ok := o.cache.Get(instanceID)
	if !ok || cachePayload == nil {
		cachePayload = make(map[string]any)
	}

	if _, exists := cachePayload[source]; !exists {
		cachePayload[source] = response
	}

	o.cache.Set(instanceID, cachePayload)
	return cachePayload, nil
}

func (o *OrchestraUsecase) handleInstanceStep(ctx context.Context, eventMsg event.GlobalEvent[any, any]) error {
	insStepExists, err := o.queries.CheckIfInstanceStepExists(ctx, eventMsg.EventID)
	if err != nil {
		return fmt.Errorf("check instance step exists: %w", err)
	}

	if !insStepExists {
		log.Println("Instance step does not exist")
	}

	eventMsgBytes, err := eventMsg.ToJSON()
	if err != nil {
		return fmt.Errorf("parse message: %w", err)
	}

	return o.queries.UpdateWorkflowInstanceStep(ctx, sqlc.UpdateWorkflowInstanceStepParams{
		Status:       eventMsg.Status,
		EventMessage: sql.NullString{String: string(eventMsgBytes), Valid: true},
		CompletedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		EventID:      eventMsg.EventID,
	})
}

func (o *OrchestraUsecase) getOrCreateWorkflowInstance(ctx context.Context, eventMsg event.GlobalEvent[any, any], wf sqlc.Workflow) (sqlc.WorkflowInstance, error) {
	if eventMsg.State == event.ORDER_CREATED.String() || eventMsg.State == event.ORDER_CANCEL.String() {
		return o.queries.CreateWorkflowInstance(ctx, sqlc.CreateWorkflowInstanceParams{
			ID:         eventMsg.InstanceID,
			WorkflowID: wf.ID,
			Status:     dto.PENDING.String(),
		})
	}

	return o.queries.FindWorkflowInstanceByID(ctx, eventMsg.InstanceID)
}

func (o *OrchestraUsecase) processSteps(ctx context.Context, eventMsg event.GlobalEvent[any, any], instance sqlc.WorkflowInstance, cachePayload map[string]any) error {
	steps, err := o.queries.FindStepsByTypeAndState(ctx, sqlc.FindStepsByTypeAndStateParams{
		Type:  eventMsg.EventType,
		State: eventMsg.State,
	})

	if err != nil {
		return fmt.Errorf("find steps: %w", err)
	}

	// handle if steps is empty, its mean all the step has done
	if len(steps) == 0 {
		err := o.processDone(ctx, eventMsg.EventType, instance.ID)

		log.Println("-> process done <-")

		if err != nil {
			return fmt.Errorf("process done: %w", err)
		}

		return nil
	}

	for _, step := range steps {
		err := o.processStep(ctx, eventMsg, instance, step, cachePayload)
		if err != nil {
			log.Printf("Error processing step %d: %v", step.StepID, err)
			continue
		}
	}

	return nil
}

func (o *OrchestraUsecase) processDone(ctx context.Context, eventType, instanceID string) error {
	var hasFailed bool

	wfiSteps, err := o.queries.FindWorkflowInstanceByTypeAndID(ctx, sqlc.FindWorkflowInstanceByTypeAndIDParams{
		Type:               eventType,
		WorkflowInstanceID: instanceID,
	})

	if err != nil {
		return fmt.Errorf("find workflow instance by type and id: %w", err)
	}

	for _, value := range wfiSteps {
		if value.InstanceStepStatus != "success" {
			hasFailed = true
			log.Printf("Step %d failed", value.StepID)
		}
	}

	if hasFailed {
		err = o.queries.UpdateWorkflowInstance(ctx, sqlc.UpdateWorkflowInstanceParams{
			Status: "failed",
			ID:     wfiSteps[0].InstanceID,
		})
	} else {
		err = o.queries.UpdateWorkflowInstance(ctx, sqlc.UpdateWorkflowInstanceParams{
			Status: "completed",
			ID:     wfiSteps[0].InstanceID,
		})
	}

	if err != nil {
		return fmt.Errorf("update workflow instance: %w", err)
	}

	return nil
}

func (o *OrchestraUsecase) processStep(ctx context.Context, eventMsg event.GlobalEvent[any, any], instance sqlc.WorkflowInstance, step sqlc.FindStepsByTypeAndStateRow, cachePayload map[string]any) error {
	keys, err := o.queries.FindPayloadKeysByStepID(ctx, step.StepID)

	if err != nil {
		return fmt.Errorf("find payload keys: %w", err)
	}

	basePayload, err := o.mergePayloads(keys, cachePayload)
	if err != nil {
		return fmt.Errorf("merge payloads: %w", err)
	}

	gevent := o.createGlobalEvent(eventMsg, basePayload, instance.ID)

	bytes, err := gevent.ToJSON()
	if err != nil {
		return fmt.Errorf("parse message: %w", err)
	}

	err = o.createWorkflowInstanceStep(ctx, gevent, step, bytes)
	if err != nil {
		return err
	}

	return o.producer.SendMessage(step.StepTopic, uuid.New().String(), bytes)
}

func (o *OrchestraUsecase) mergePayloads(keys []string, cachePayload map[string]any) (any, error) {
	var basePayloads []any
	for _, key := range keys {
		basePayloads = append(basePayloads, cachePayload[key])
	}

	var basePayload any
	err := pkg.MergeJSON(&basePayload, basePayloads...)
	return basePayload, err
}

func (o *OrchestraUsecase) createGlobalEvent(eventMsg event.GlobalEvent[any, any], basePayload any, instanceID string) event.GlobalEvent[any, any] {
	gevent := event.NewGlobalEvent[any, any](
		"redirect",
		eventMsg.Status,
		event.BasePayload[any, any]{
			Request: basePayload,
		},
	)

	gevent.State = eventMsg.State
	gevent.EventType = eventMsg.EventType
	gevent.StatusCode = eventMsg.StatusCode
	gevent.InstanceID = instanceID

	return gevent
}

func (o *OrchestraUsecase) createWorkflowInstanceStep(ctx context.Context, gevent event.GlobalEvent[any, any], step sqlc.FindStepsByTypeAndStateRow, eventMessage []byte) error {
	_, err := o.queries.CreateWorkflowInstanceStep(ctx, sqlc.CreateWorkflowInstanceStepParams{
		WorkflowInstanceID: gevent.InstanceID,
		EventID:            gevent.EventID,
		Status:             dto.IN_PROGRESS.String(),
		StepID:             step.StepID,
		EventMessage:       sql.NullString{String: string(eventMessage), Valid: true},
		StartedAt:          sql.NullTime{Time: time.Now(), Valid: true},
	})
	return err
}
