package dto

import "orchestra-svc/internal/repository/sqlc"

type WorkflowResponse struct {
	Workflow sqlc.Workflow `json:"workflow"`
	Steps    []sqlc.Step   `json:"steps"`
}
