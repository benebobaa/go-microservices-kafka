// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"context"
)

type Querier interface {
	CheckIfInstanceStepExists(ctx context.Context, eventID string) (bool, error)
	CreateProcessLog(ctx context.Context, arg CreateProcessLogParams) error
	CreateWorkflowInstance(ctx context.Context, arg CreateWorkflowInstanceParams) (WorkflowInstance, error)
	CreateWorkflowInstanceStep(ctx context.Context, arg CreateWorkflowInstanceStepParams) (WorkflowInstanceStep, error)
	FindInstanceStepByEventID(ctx context.Context, eventID string) (WorkflowInstanceStep, error)
	FindInstanceStepByID(ctx context.Context, workflowInstanceID string) ([]WorkflowInstanceStep, error)
	FindPayloadKeysByStepID(ctx context.Context, stepID int32) ([]string, error)
	FindStepsByTypeAndState(ctx context.Context, arg FindStepsByTypeAndStateParams) ([]FindStepsByTypeAndStateRow, error)
	FindWorkflowByType(ctx context.Context, type_ string) (Workflow, error)
	FindWorkflowInstanceByID(ctx context.Context, id string) (WorkflowInstance, error)
	FindWorkflowInstanceByTypeAndID(ctx context.Context, arg FindWorkflowInstanceByTypeAndIDParams) ([]FindWorkflowInstanceByTypeAndIDRow, error)
	FindWorkflowInstanceStepsByEventIDAndInsID(ctx context.Context, arg FindWorkflowInstanceStepsByEventIDAndInsIDParams) (FindWorkflowInstanceStepsByEventIDAndInsIDRow, error)
	UpdateWorkflowInstance(ctx context.Context, arg UpdateWorkflowInstanceParams) error
	UpdateWorkflowInstanceStep(ctx context.Context, arg UpdateWorkflowInstanceStepParams) error
}

var _ Querier = (*Queries)(nil)
