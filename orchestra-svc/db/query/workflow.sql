
-- name: FindWorkflowByType :one
SELECT * FROM workflows
WHERE type = $1 LIMIT 1;
