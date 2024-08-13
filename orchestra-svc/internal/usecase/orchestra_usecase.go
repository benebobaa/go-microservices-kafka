package usecase

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"log"
	"orchestra-svc/internal/dto"
	"orchestra-svc/internal/dto/event"
	"orchestra-svc/internal/repository/sqlc"
	"orchestra-svc/pkg/producer"
)

type OrchestraUsecase struct {
	queries             sqlc.Querier
	userProductProducer *producer.KafkaProducer
}

func NewOrchestraUsecase(
	q sqlc.Querier,
	upp *producer.KafkaProducer,
) *OrchestraUsecase {
	return &OrchestraUsecase{
		queries:             q,
		userProductProducer: upp,
	}
}

func (o *OrchestraUsecase) ProcessWorkflow(ctx context.Context, eventMsg event.GlobalEvent[any]) error {

	var instance sqlc.WorkflowInstance
	var err error

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
		eventMsg.InstanceID = instance.ID
		bytes, err := eventMsg.ToJSON()

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
