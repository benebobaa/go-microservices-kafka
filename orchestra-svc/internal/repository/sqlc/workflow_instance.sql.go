// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: workflow_instance.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createWorkflowInstance = `-- name: CreateWorkflowInstance :one
INSERT INTO workflow_instances (id, workflow_id, status)
VALUES
    ($1, $2, $3) RETURNING id, workflow_id, status, created_at, updated_at
`

type CreateWorkflowInstanceParams struct {
	ID         string `json:"id"`
	WorkflowID int32  `json:"workflow_id"`
	Status     string `json:"status"`
}

func (q *Queries) CreateWorkflowInstance(ctx context.Context, arg CreateWorkflowInstanceParams) (WorkflowInstance, error) {
	row := q.db.QueryRowContext(ctx, createWorkflowInstance, arg.ID, arg.WorkflowID, arg.Status)
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
INSERT INTO workflow_instance_steps (workflow_instance_id, step_id, status, event_message)
VALUES
    ($1, $2, $3, $4) RETURNING id, workflow_instance_id, step_id, status, event_message, created_at
`

type CreateWorkflowInstanceStepParams struct {
	WorkflowInstanceID string         `json:"workflow_instance_id"`
	StepID             int32          `json:"step_id"`
	Status             string         `json:"status"`
	EventMessage       sql.NullString `json:"event_message"`
}

func (q *Queries) CreateWorkflowInstanceStep(ctx context.Context, arg CreateWorkflowInstanceStepParams) (WorkflowInstanceStep, error) {
	row := q.db.QueryRowContext(ctx, createWorkflowInstanceStep,
		arg.WorkflowInstanceID,
		arg.StepID,
		arg.Status,
		arg.EventMessage,
	)
	var i WorkflowInstanceStep
	err := row.Scan(
		&i.ID,
		&i.WorkflowInstanceID,
		&i.StepID,
		&i.Status,
		&i.EventMessage,
		&i.CreatedAt,
	)
	return i, err
}

const findInstanceStepByID = `-- name: FindInstanceStepByID :many
SELECT id, workflow_instance_id, step_id, status, event_message, created_at FROM workflow_instance_steps
WHERE workflow_instance_id = $1
`

func (q *Queries) FindInstanceStepByID(ctx context.Context, workflowInstanceID string) ([]WorkflowInstanceStep, error) {
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
			&i.StepID,
			&i.Status,
			&i.EventMessage,
			&i.CreatedAt,
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

const findStepsByState = `-- name: FindStepsByState :many
SELECT DISTINCT
    sa.state,
    s.id AS step_id,
    s.service AS service,
    s.name AS step_name,
    s.description AS step_description,
    s.topic AS step_topic
FROM
    state_actions sa
        JOIN steps s ON sa.step_id = s.id
WHERE
    sa.state = $1
ORDER BY
    sa.state
`

type FindStepsByStateRow struct {
	State           string `json:"state"`
	StepID          int32  `json:"step_id"`
	Service         string `json:"service"`
	StepName        string `json:"step_name"`
	StepDescription string `json:"step_description"`
	StepTopic       string `json:"step_topic"`
}

func (q *Queries) FindStepsByState(ctx context.Context, state string) ([]FindStepsByStateRow, error) {
	rows, err := q.db.QueryContext(ctx, findStepsByState, state)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FindStepsByStateRow{}
	for rows.Next() {
		var i FindStepsByStateRow
		if err := rows.Scan(
			&i.State,
			&i.StepID,
			&i.Service,
			&i.StepName,
			&i.StepDescription,
			&i.StepTopic,
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

const findWorkflowInstanceByID = `-- name: FindWorkflowInstanceByID :one
SELECT id, workflow_id, status, created_at, updated_at FROM workflow_instances
WHERE id = $1 LIMIT 1
`

func (q *Queries) FindWorkflowInstanceByID(ctx context.Context, id string) (WorkflowInstance, error) {
	row := q.db.QueryRowContext(ctx, findWorkflowInstanceByID, id)
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
