-- name: CreateWorkflowType :one
INSERT INTO workflow_types (name) 
VALUES 
    ($1) RETURNING *;

-- name: CreateWorkflow :one
INSERT INTO workflows (workflow_type_id, name, description) 
VALUES 
    ($1, $2, $3) RETURNING *;

-- name: CreateStep :one
INSERT INTO steps (name, description) 
VALUES 
    ($1, $2) RETURNING *;

-- name: CreateWorkflowStep :one
INSERT INTO workflow_steps (workflow_id, step_id, order_index)
VALUES 
    ($1, $2, $3) RETURNING *;

-- name: GetWorkflowStepByType :many
SELECT
    wt.id AS workflow_type_id,
    ws.order_index AS step_order,
    s.id AS step_id,
    s.name AS step_name,
    s.description AS step_description
FROM
    workflow_types wt
JOIN
    workflows w ON wt.id = w.workflow_type_id
JOIN
    workflow_steps ws ON w.id = ws.workflow_id
JOIN
    steps s ON ws.step_id = s.id
WHERE
    wt.name = $1
ORDER BY
    ws.order_index;
