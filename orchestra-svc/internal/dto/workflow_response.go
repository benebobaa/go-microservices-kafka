package dto

import "orchestra-svc/internal/repository/sqlc"

type WorkflowResponse struct {
	Type     sqlc.WorkflowType `json:"type"`
	Workflow sqlc.Workflow     `json:"workflow"`
	Steps    []sqlc.Step       `json:"steps"`
}
