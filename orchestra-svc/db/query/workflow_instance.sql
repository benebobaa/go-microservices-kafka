-- name: CreateWorkflowInstance :one
INSERT INTO workflow_instances (workflow_id, status) 
VALUES 
    ($1, $2) RETURNING *;

-- name: CreateWorkflowInstanceStep :one
INSERT INTO workflow_instance_steps (workflow_instance_id, workflow_step_id, status) 
VALUES 
    ($1, $2, $3) RETURNING *;

-- name: FindInstanceStepByID :many
SELECT * FROM workflow_instance_steps
WHERE workflow_instance_id = $1;
