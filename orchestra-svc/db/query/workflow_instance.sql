
-- name: CreateWorkflowInstance :one
INSERT INTO workflow_instances (id, workflow_id, status)
VALUES
    ($1, $2, $3) RETURNING *;

-- name: CreateWorkflowInstanceStep :one
INSERT INTO workflow_instance_steps (workflow_instance_id,event_id, step_id, status, event_message, started_at, completed_at)
VALUES
    ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: CheckIfInstanceStepExists :one
SELECT EXISTS(SELECT 1 FROM workflow_instance_steps WHERE event_id = $1) AS exists;

-- name: UpdateWorkflowInstanceStep :exec
UPDATE workflow_instance_steps
SET
    status = $1,
    event_message = $2,
    completed_at = $3
WHERE
    event_id = $4;

-- name: FindInstanceStepByID :many
SELECT * FROM workflow_instance_steps
WHERE workflow_instance_id = $1;

-- name: FindWorkflowInstanceByID :one
SELECT * FROM workflow_instances
WHERE id = $1 LIMIT 1;
