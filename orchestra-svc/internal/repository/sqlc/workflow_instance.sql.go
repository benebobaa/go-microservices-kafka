// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: workflow_instance.sql

package sqlc

import (
	"context"
	"database/sql"
)

const checkIfInstanceStepExists = `-- name: CheckIfInstanceStepExists :one
SELECT EXISTS(SELECT 1 FROM workflow_instance_steps WHERE event_id = $1) AS exists
`

func (q *Queries) CheckIfInstanceStepExists(ctx context.Context, eventID string) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkIfInstanceStepExists, eventID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

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
INSERT INTO workflow_instance_steps (workflow_instance_id,event_id, step_id, status, event_message, started_at, completed_at)
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING id, event_id, status_code, response, workflow_instance_id, step_id, status, event_message, started_at, completed_at
`

type CreateWorkflowInstanceStepParams struct {
	WorkflowInstanceID string         `json:"workflow_instance_id"`
	EventID            string         `json:"event_id"`
	StepID             int32          `json:"step_id"`
	Status             string         `json:"status"`
	EventMessage       sql.NullString `json:"event_message"`
	StartedAt          sql.NullTime   `json:"started_at"`
	CompletedAt        sql.NullTime   `json:"completed_at"`
}

func (q *Queries) CreateWorkflowInstanceStep(ctx context.Context, arg CreateWorkflowInstanceStepParams) (WorkflowInstanceStep, error) {
	row := q.db.QueryRowContext(ctx, createWorkflowInstanceStep,
		arg.WorkflowInstanceID,
		arg.EventID,
		arg.StepID,
		arg.Status,
		arg.EventMessage,
		arg.StartedAt,
		arg.CompletedAt,
	)
	var i WorkflowInstanceStep
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.StatusCode,
		&i.Response,
		&i.WorkflowInstanceID,
		&i.StepID,
		&i.Status,
		&i.EventMessage,
		&i.StartedAt,
		&i.CompletedAt,
	)
	return i, err
}

const findInstanceStepByEventID = `-- name: FindInstanceStepByEventID :one
SELECT id, event_id, status_code, response, workflow_instance_id, step_id, status, event_message, started_at, completed_at FROM workflow_instance_steps
WHERE event_id = $1 LIMIT 1
`

func (q *Queries) FindInstanceStepByEventID(ctx context.Context, eventID string) (WorkflowInstanceStep, error) {
	row := q.db.QueryRowContext(ctx, findInstanceStepByEventID, eventID)
	var i WorkflowInstanceStep
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.StatusCode,
		&i.Response,
		&i.WorkflowInstanceID,
		&i.StepID,
		&i.Status,
		&i.EventMessage,
		&i.StartedAt,
		&i.CompletedAt,
	)
	return i, err
}

const findInstanceStepByID = `-- name: FindInstanceStepByID :many
SELECT id, event_id, status_code, response, workflow_instance_id, step_id, status, event_message, started_at, completed_at FROM workflow_instance_steps
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
			&i.EventID,
			&i.StatusCode,
			&i.Response,
			&i.WorkflowInstanceID,
			&i.StepID,
			&i.Status,
			&i.EventMessage,
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

const findWorkflowInstanceByTypeAndID = `-- name: FindWorkflowInstanceByTypeAndID :many
SELECT
    w.id as workflow_id,
    w.type as workflow_type,
    wi.id as instance_id,
    wi.status as instance_status,
    wis.status as instance_step_status,
    wis.step_id as step_id
FROM workflows w
JOIN workflow_instances wi on w.id = wi.workflow_id
JOIN workflow_instance_steps wis ON wi.id = wis.workflow_instance_id
WHERE w.type = $1 AND wis.workflow_instance_id = $2
`

type FindWorkflowInstanceByTypeAndIDParams struct {
	Type               string `json:"type"`
	WorkflowInstanceID string `json:"workflow_instance_id"`
}

type FindWorkflowInstanceByTypeAndIDRow struct {
	WorkflowID         int32  `json:"workflow_id"`
	WorkflowType       string `json:"workflow_type"`
	InstanceID         string `json:"instance_id"`
	InstanceStatus     string `json:"instance_status"`
	InstanceStepStatus string `json:"instance_step_status"`
	StepID             int32  `json:"step_id"`
}

func (q *Queries) FindWorkflowInstanceByTypeAndID(ctx context.Context, arg FindWorkflowInstanceByTypeAndIDParams) ([]FindWorkflowInstanceByTypeAndIDRow, error) {
	rows, err := q.db.QueryContext(ctx, findWorkflowInstanceByTypeAndID, arg.Type, arg.WorkflowInstanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FindWorkflowInstanceByTypeAndIDRow{}
	for rows.Next() {
		var i FindWorkflowInstanceByTypeAndIDRow
		if err := rows.Scan(
			&i.WorkflowID,
			&i.WorkflowType,
			&i.InstanceID,
			&i.InstanceStatus,
			&i.InstanceStepStatus,
			&i.StepID,
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

const updateWorkflowInstance = `-- name: UpdateWorkflowInstance :exec
UPDATE workflow_instances
SET
    status = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $2
`

type UpdateWorkflowInstanceParams struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

func (q *Queries) UpdateWorkflowInstance(ctx context.Context, arg UpdateWorkflowInstanceParams) error {
	_, err := q.db.ExecContext(ctx, updateWorkflowInstance, arg.Status, arg.ID)
	return err
}

const updateWorkflowInstanceStep = `-- name: UpdateWorkflowInstanceStep :exec
UPDATE workflow_instance_steps
SET
    status = $1,
    event_message = $2,
    status_code = $3,
    response = $4,
    completed_at = $5
WHERE
    event_id = $6
`

type UpdateWorkflowInstanceStepParams struct {
	Status       string         `json:"status"`
	EventMessage sql.NullString `json:"event_message"`
	StatusCode   sql.NullInt32  `json:"status_code"`
	Response     sql.NullString `json:"response"`
	CompletedAt  sql.NullTime   `json:"completed_at"`
	EventID      string         `json:"event_id"`
}

func (q *Queries) UpdateWorkflowInstanceStep(ctx context.Context, arg UpdateWorkflowInstanceStepParams) error {
	_, err := q.db.ExecContext(ctx, updateWorkflowInstanceStep,
		arg.Status,
		arg.EventMessage,
		arg.StatusCode,
		arg.Response,
		arg.CompletedAt,
		arg.EventID,
	)
	return err
}
