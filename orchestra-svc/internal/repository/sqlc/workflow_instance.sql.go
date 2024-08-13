// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: workflow_instance.sql

package sqlc

import (
	"context"
)

const createWorkflowInstance = `-- name: CreateWorkflowInstance :one
INSERT INTO workflow_instances (workflow_id, status) 
VALUES 
    ($1, $2) RETURNING id, workflow_id, status, created_at, updated_at
`

type CreateWorkflowInstanceParams struct {
	WorkflowID int32  `json:"workflow_id"`
	Status     string `json:"status"`
}

func (q *Queries) CreateWorkflowInstance(ctx context.Context, arg CreateWorkflowInstanceParams) (WorkflowInstance, error) {
	row := q.db.QueryRowContext(ctx, createWorkflowInstance, arg.WorkflowID, arg.Status)
	var i WorkflowInstance
	err := row.Scan(
		&i.ID,
		&i.WorkflowID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createWorkflowInstanceStep = `-- name: CreateWorkflowInstanceStep :one
INSERT INTO workflow_instance_steps (workflow_instance_id, workflow_step_id, status) 
VALUES 
    ($1, $2, $3) RETURNING id, workflow_instance_id, workflow_step_id, status, started_at, completed_at
`

type CreateWorkflowInstanceStepParams struct {
	WorkflowInstanceID int32  `json:"workflow_instance_id"`
	WorkflowStepID     int32  `json:"workflow_step_id"`
	Status             string `json:"status"`
}

func (q *Queries) CreateWorkflowInstanceStep(ctx context.Context, arg CreateWorkflowInstanceStepParams) (WorkflowInstanceStep, error) {
	row := q.db.QueryRowContext(ctx, createWorkflowInstanceStep, arg.WorkflowInstanceID, arg.WorkflowStepID, arg.Status)
	var i WorkflowInstanceStep
	err := row.Scan(
		&i.ID,
		&i.WorkflowInstanceID,
		&i.WorkflowStepID,
		&i.Status,
		&i.StartedAt,
		&i.CompletedAt,
	)
	return i, err
}

const findInstanceStepByID = `-- name: FindInstanceStepByID :many
SELECT id, workflow_instance_id, workflow_step_id, status, started_at, completed_at FROM workflow_instance_steps
WHERE workflow_instance_id = $1
`

func (q *Queries) FindInstanceStepByID(ctx context.Context, workflowInstanceID int32) ([]WorkflowInstanceStep, error) {
	rows, err := q.db.QueryContext(ctx, findInstanceStepByID, workflowInstanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []WorkflowInstanceStep{}
	for rows.Next() {
		var i WorkflowInstanceStep
		if err := rows.Scan(
			&i.ID,
			&i.WorkflowInstanceID,
			&i.WorkflowStepID,
			&i.Status,
			&i.StartedAt,
			&i.CompletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
