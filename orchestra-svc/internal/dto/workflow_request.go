package dto

type WorkflowRequest struct {
	Type        string `json:"name" valo:"notblank"`
	Description string `json:"description" valo:"notblank"`
	Step        []Step `json:"step" valo:"sizeMin=1,valid"`
}

type Step struct {
	Name        string `json:"name" valo:"notblank"`
	State       string `json:"state" valo:"notblank"`
	Description string `json:"description" valo:"notblank"`
}
