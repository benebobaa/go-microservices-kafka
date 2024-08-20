package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
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
	queries  sqlc.Store
	producer *producer.KafkaProducer
	cache    *cache.PayloadCacher
}

func NewOrchestraUsecase(q sqlc.Store, p *producer.KafkaProducer, c *cache.PayloadCacher) *OrchestraUsecase {
	return &OrchestraUsecase{
		queries:  q,
		producer: p,
		cache:    c,
	}
}

func (o *OrchestraUsecase) ProcessWorkflow(ctx context.Context, eventMsg event.GlobalEvent[any, any]) error {

	log.Println("Processing workflow: ", eventMsg.EventType)
	log.Println("Processing state: ", eventMsg.State)

	err := o.logDB(ctx, eventMsg)

	if err != nil {
		log.Println("Error logging to db: ", err)
	}

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

	cachePayload[source] = response

	o.cache.Set(instanceID, cachePayload)
	return cachePayload, nil
}

func (o *OrchestraUsecase) handleInstanceStep(ctx context.Context, eventMsg event.GlobalEvent[any, any]) error {
	instanceStep, err := o.queries.FindInstanceStepByEventID(ctx, eventMsg.EventID)
	if err != nil {
		log.Println("Error find instance step by event id: ", err)
		//return fmt.Errorf("check instance step exists: %w", err)
		return err
	}

	if eventMsg.StatusCode >= 500 {
		log.Println("err server 500: ", instanceStep.StepID)
	}

	eventMsgBytes, err := eventMsg.ToJSON()
	if err != nil {
		return fmt.Errorf("parse message: %w", err)
	}

	responseMsg, err := json.Marshal(eventMsg.Payload.Response)

	if err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	log.Println("step id: ", instanceStep.StepID)
	log.Println("process time: ", time.Since(instanceStep.StartedAt.Time))

	return o.queries.UpdateWorkflowInstanceStep(ctx, sqlc.UpdateWorkflowInstanceStepParams{
		Status:       eventMsg.Status,
		StatusCode:   sql.NullInt32{Int32: int32(eventMsg.StatusCode), Valid: true},
		Response:     sql.NullString{String: string(responseMsg), Valid: true},
		EventMessage: sql.NullString{String: string(eventMsgBytes), Valid: true},
		StartedAt:    instanceStep.StartedAt,
		CompletedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		EventID:      eventMsg.EventID,
	})
}

func (o *OrchestraUsecase) getOrCreateWorkflowInstance(ctx context.Context, eventMsg event.GlobalEvent[any, any], wf sqlc.Workflow) (sqlc.WorkflowInstance, error) {
	if eventMsg.State == event.ORDER_CREATED.String() || eventMsg.State == event.ORDER_CANCEL.String() || eventMsg.State == event.BANK_REGIS_CREATED.String() {
		return o.queries.CreateWorkflowInstance(ctx, sqlc.CreateWorkflowInstanceParams{
			ID:         eventMsg.InstanceID,
			WorkflowID: wf.ID,
			Status:     dto.IN_PROGRESS.String(),
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

func (o *OrchestraUsecase) logDB(ctx context.Context, globalEvent event.GlobalEvent[any, any]) error {
	bytes, err := json.Marshal(globalEvent)

	if err != nil {
		return fmt.Errorf("parse message: %w", err)
	}

	return o.queries.CreateProcessLog(ctx, sqlc.CreateProcessLogParams{
		EventID:            globalEvent.EventID,
		WorkflowInstanceID: globalEvent.InstanceID,
		State:              globalEvent.State,
		StatusCode: sql.NullInt32{
			Int32: int32(globalEvent.StatusCode),
			Valid: true,
		},
		Status:       globalEvent.Status,
		EventMessage: string(bytes),
	})
}
