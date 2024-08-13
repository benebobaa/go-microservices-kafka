package usecase

import (
	"context"
	"errors"
	"log"
	"orchestra-svc/internal/dto"
	"orchestra-svc/internal/dto/event"
	"orchestra-svc/internal/repository/sqlc"
	"orchestra-svc/pkg/producer"
)

type OrchestraUsecase struct {
	queries      sqlc.Querier
	userProducer *producer.KafkaProducer
}

func NewOrchestraUsecase(
	q sqlc.Querier,
	up *producer.KafkaProducer,
) *OrchestraUsecase {
	return &OrchestraUsecase{
		queries:      q,
		userProducer: up,
	}
}

func (o *OrchestraUsecase) ProcessWorkflow(ctx context.Context, eventMsg event.GlobalEvent[any]) error {

	steps, err := o.queries.FindInstanceStepByID(ctx, eventMsg.InstanceID)

	if err != nil {
		return err
	}

	if len(steps) < 1 {
		return errors.New("workflow steps is empty")
	}

	for _, value := range steps {
		log.Println("step value: ", value)

	}

	return nil
}

func (o *OrchestraUsecase) ProcessNewWorkflow(ctx context.Context, eventMsg event.GlobalEvent[any]) error {

	steps, err := o.queries.GetWorkflowStepByType(ctx, eventMsg.EventType)

	if len(steps) < 1 {
		return errors.New("workflow steps is empty")
	}

	if err != nil {
		return err
	}

	wfi, err := o.queries.CreateWorkflowInstance(ctx, sqlc.CreateWorkflowInstanceParams{
		WorkflowID: steps[0].WorkflowTypeID,
		Status:     dto.PENDING.String(),
	})

	if err != nil {
		return err
	}

	eventMsg.InstanceID = wfi.ID

	for _, value := range steps {

		_, err = o.queries.CreateWorkflowInstanceStep(ctx, sqlc.CreateWorkflowInstanceStepParams{
			WorkflowInstanceID: wfi.ID,
			Status:             dto.PENDING.String(),
			WorkflowStepID:     value.StepID,
		})

		if err != nil {
			log.Println("err CreateWorkflowInstanceStep: ", err.Error())
		}
	}

	err = o.ProcessWorkflow(ctx, eventMsg)

	if err != nil {
		return err
	}

	return nil
}
