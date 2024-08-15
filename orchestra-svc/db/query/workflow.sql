
-- name: FindWorkflowByType :one
SELECT * FROM workflows
WHERE type = $1 LIMIT 1;

-- name: FindPayloadKeysByStepID :many
SELECT key FROM payload_keys WHERE step_id = $1;

-- name: FindStepsByTypeAndState :many
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
    sa.state;

