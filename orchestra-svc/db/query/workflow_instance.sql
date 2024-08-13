
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

-- name: FindStepsByState :many
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
    sa.state;
