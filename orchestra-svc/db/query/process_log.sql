-- name: CreateProcessLog :exec
INSERT INTO process_logs (
    event_id,
    workflow_instance_id,
    state,
    status_code,
    status,
    event_message
) VALUES ($1, $2, $3, $4, $5, $6 );
