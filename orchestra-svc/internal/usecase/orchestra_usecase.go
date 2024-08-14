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

	if eventMsg.State == event.ORDER_CREATED.String() {
		instance, err = o.queries.CreateWorkflowInstance(ctx, sqlc.CreateWorkflowInstanceParams{
			ID:         eventMsg.InstanceID,
			WorkflowID: wf.ID,
			Status:     dto.PENDING.String(),
		})

		if err != nil {
			log.Println("Error when create workflow instance: ", err.Error())
			return err
		}
	}

	if eventMsg.InstanceID != "" {
		instance, err = o.queries.FindWorkflowInstanceByID(ctx, eventMsg.InstanceID)
	}

	if err != nil {
		log.Println("Error when find workflow instance: ", err.Error())
		return err
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

		log.Println("step: ", step)
		log.Println("keys: ", keys)

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
			"success",
			event.BasePayload[any, any]{
				Request: basePayload,
			},
		)

		gevent.State = eventMsg.State
		gevent.EventType = eventMsg.EventType
		gevent.StatusCode = 200
		gevent.InstanceID = instance.ID
		bytes, err := gevent.ToJSON()

		if err != nil {
			log.Println("Error when parse message: ", err.Error())
			continue
		}

		_, err = o.queries.CreateWorkflowInstanceStep(ctx, sqlc.CreateWorkflowInstanceStepParams{
			WorkflowInstanceID: eventMsg.InstanceID,
			Status:             dto.IN_PROGRESS.String(),
			StepID:             step.StepID,
			EventMessage: sql.NullString{
				String: string(bytes),
				Valid:  true,
			},
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
