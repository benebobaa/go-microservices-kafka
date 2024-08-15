package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"orchestra-svc/internal/dto"
	"orchestra-svc/internal/dto/event"
	"orchestra-svc/internal/repository/cache"
	"orchestra-svc/internal/repository/sqlc"
	"orchestra-svc/pkg/producer"
	"time"
)

type OrchestraUsecase struct {
	queries             sqlc.Querier
	userProductProducer *producer.KafkaProducer
	cache               *cache.PayloadCache
}

func NewOrchestraUsecase(
	q sqlc.Querier,
	upp *producer.KafkaProducer,
	c *cache.PayloadCache,
) *OrchestraUsecase {
	return &OrchestraUsecase{
		queries:             q,
		userProductProducer: upp,
		cache:               c,
	}
}

func (o *OrchestraUsecase) ProcessWorkflow(ctx context.Context, eventMsg event.GlobalEvent[any, any]) error {

	var instance sqlc.WorkflowInstance
	var err error

	cachePayload, ok := o.cache.Get(eventMsg.InstanceID)

	if !ok || cachePayload == nil {
		cachePayload = make(map[string]any)
	}

	_, exists := cachePayload[eventMsg.Source]
	if !exists {
		cachePayload[eventMsg.Source] = eventMsg.Payload.Response
	}

	o.cache.Set(eventMsg.InstanceID, cachePayload)

	wf, err := o.queries.FindWorkflowByType(ctx, eventMsg.EventType)

	if err != nil {
		log.Println("Error when find workflow: ", err.Error())
		return err
	}

	eventMsgBytes, err := eventMsg.ToJSON()
	if err != nil {
		log.Println("Error when parse message: ", err.Error())
		return err
	}

	stepExists, err := o.queries.CheckIfInstanceStepExists(ctx, eventMsg.EventID)
	if err != nil {
		log.Println("Error when check instance step exists: ", err.Error())
		return err
	}

	if !stepExists {
		log.Println("-> step not exists <-")
	} else {
		err = o.queries.UpdateWorkflowInstanceStep(ctx, sqlc.UpdateWorkflowInstanceStepParams{
			Status: eventMsg.Status,
			EventMessage: sql.NullString{
				String: string(eventMsgBytes),
				Valid:  true,
			},
			CompletedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			EventID: eventMsg.EventID,
		})

		if err != nil {
			log.Println("Error when update instance step: ", err.Error())
			return err
		}
	}

	if eventMsg.State == event.ORDER_CREATED.String() {
		// this mean init process order process

		instance, err = o.queries.CreateWorkflowInstance(ctx, sqlc.CreateWorkflowInstanceParams{
			ID:         eventMsg.InstanceID,
			WorkflowID: wf.ID,
			Status:     dto.PENDING.String(),
		})

		if err != nil {
			log.Println("Error when create workflow instance: ", err.Error())
			return err
		}
	} else if eventMsg.State == event.ORDER_CANCELLED.String() {
		// this mean init process for order cancel

		instance, err = o.queries.CreateWorkflowInstance(ctx, sqlc.CreateWorkflowInstanceParams{
			ID:         eventMsg.InstanceID,
			WorkflowID: wf.ID,
			Status:     dto.PENDING.String(),
		})

		if err != nil {
			log.Println("Error when create workflow instance: ", err.Error())
			return err
		}
	} else {

		instance, err = o.queries.FindWorkflowInstanceByID(ctx, eventMsg.InstanceID)

		if err != nil {
			log.Println("Error when find workflow instance: ", err.Error())
			return err
		}
	}

	steps, err := o.queries.FindStepsByState(ctx, eventMsg.State)

	log.Println("steps: ", steps)

	if err != nil {
		log.Println("Error when find steps: ", err.Error())
		return err
	}

	for _, step := range steps {
		var basePayloads []any
		var basePayload any
		keys, err := o.queries.FindPayloadKeysByStepID(ctx, step.StepID)

		payloadCache, _ := o.cache.Get(eventMsg.InstanceID)

		for _, key := range keys {
			basePayloads = append(basePayloads, payloadCache[key])
		}

		err = mergeJSON(&basePayload, basePayloads...)

		if err != nil {
			log.Println("Error when merge JSON: ", err.Error())
			return err
		}

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
		gevent.InstanceID = instance.ID
		bytes, err := gevent.ToJSON()

		if err != nil {
			log.Println("Error when parse message: ", err.Error())
			continue
		}

		_, err = o.queries.CreateWorkflowInstanceStep(ctx, sqlc.CreateWorkflowInstanceStepParams{
			WorkflowInstanceID: gevent.InstanceID,
			EventID:            gevent.EventID,
			Status:             dto.IN_PROGRESS.String(),
			StepID:             step.StepID,
			EventMessage: sql.NullString{
				String: string(bytes),
				Valid:  true,
			},
			StartedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})

		if err != nil {
			log.Println("Error when create instance step: ", err.Error())
			continue
		}

		err = o.userProductProducer.SendMessage(
			step.StepTopic,
			uuid.New().String(),
			bytes,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func mergeJSON(dst interface{}, sources ...interface{}) error {
	mergedMap := make(map[string]interface{})

	for _, src := range sources {
		jsonData, err := json.Marshal(src)
		if err != nil {
			return err
		}

		tempMap := make(map[string]interface{})
		if err := json.Unmarshal(jsonData, &tempMap); err != nil {
			return err
		}

		for key, value := range tempMap {
			mergedMap[key] = value
		}
	}

	finalJSON, err := json.Marshal(mergedMap)
	if err != nil {
		return err
	}

	return json.Unmarshal(finalJSON, dst)
}
