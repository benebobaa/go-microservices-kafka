// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: workflow.sql

package sqlc

import (
	"context"
)

const findPayloadKeysByStepID = `-- name: FindPayloadKeysByStepID :many
SELECT key FROM payload_keys WHERE step_id = $1
`

func (q *Queries) FindPayloadKeysByStepID(ctx context.Context, stepID int32) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, findPayloadKeysByStepID, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		items = append(items, key)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findStepsByTypeAndState = `-- name: FindStepsByTypeAndState :many
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
    sa.type = $1 AND
    sa.state = $2
ORDER BY
    sa.state
`

type FindStepsByTypeAndStateParams struct {
	Type  string `json:"type"`
	State string `json:"state"`
}

type FindStepsByTypeAndStateRow struct {
	State           string `json:"state"`
	StepID          int32  `json:"step_id"`
	Service         string `json:"service"`
	StepName        string `json:"step_name"`
	StepDescription string `json:"step_description"`
	StepTopic       string `json:"step_topic"`
}

func (q *Queries) FindStepsByTypeAndState(ctx context.Context, arg FindStepsByTypeAndStateParams) ([]FindStepsByTypeAndStateRow, error) {
	rows, err := q.db.QueryContext(ctx, findStepsByTypeAndState, arg.Type, arg.State)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FindStepsByTypeAndStateRow{}
	for rows.Next() {
		var i FindStepsByTypeAndStateRow
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

const findWorkflowByType = `-- name: FindWorkflowByType :one
SELECT id, type, description, created_at, updated_at FROM workflows
WHERE type = $1 LIMIT 1
`

func (q *Queries) FindWorkflowByType(ctx context.Context, type_ string) (Workflow, error) {
	row := q.db.QueryRowContext(ctx, findWorkflowByType, type_)
	var i Workflow
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
