
-- name: CreateWorkflowInstance :one
INSERT INTO workflow_instances (id, workflow_id, status)
VALUES
    ($1, $2, $3) RETURNING *;

-- name: CreateWorkflowInstanceStep :one
INSERT INTO workflow_instance_steps (workflow_instance_id, step_id, status, event_message)
VALUES
    ($1, $2, $3, $4) RETURNING *;

-- name: FindInstanceStepByID :many
SELECT * FROM workflow_instance_steps
WHERE workflow_instance_id = $1;

-- name: FindWorkflowInstanceByID :one
SELECT * FROM workflow_instances
WHERE id = $1 LIMIT 1;
