package usecase

import (
	"orchestra-svc/internal/dto"
	"orchestra-svc/internal/repository/sqlc"

	"golang.org/x/net/context"
)

type WorkflowUsecase struct {
	queries sqlc.Querier
}

func NewWorkflowUsecase(queries sqlc.Querier) *WorkflowUsecase {
	return &WorkflowUsecase{queries: queries}
}

func (w *WorkflowUsecase) CreateWorkflow(ctx context.Context, workflow *dto.WorkflowRequest) (*dto.WorkflowResponse, error) {

	//wt, err := w.queries.CreateWorkflowType(ctx, workflow.Name)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//wf, err := w.queries.CreateWorkflow(ctx, sqlc.CreateWorkflowParams{
	//	WorkflowTypeID: wt.ID,
	//	Name:           workflow.Workflow.Name,
	//	Description:    sql.NullString{String: workflow.Workflow.Description, Valid: true},
	//})
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//var steps []sqlc.Step
	//
	//for i, v := range workflow.Step {
	//
	//	ws, err := w.queries.CreateStep(ctx, sqlc.CreateStepParams{
	//		Name:        v.Name,
	//		Description: sql.NullString{String: v.Description, Valid: true},
	//	})
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	steps = append(steps, ws)
	//
	//	_, err = w.queries.CreateWorkflowStep(ctx, sqlc.CreateWorkflowStepParams{
	//		WorkflowID: wf.ID,
	//		StepID:     ws.ID,
	//		OrderIndex: int32(i) + 1,
	//	})
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	//return &dto.WorkflowResponse{
	//	Type:     wt,
	//	Workflow: wf,
	//	Steps:    steps,
	//}, nil

	return nil, nil
}

//func (w *WorkflowUsecase) GetByType(ctx context.Context, name string) ([]sqlc.FindWorkflowStepByTypeRow, error) {
//	return nil, nil
//}
