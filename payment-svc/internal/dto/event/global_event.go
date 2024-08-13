package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type GlobalEvent[T any] struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
	Action    string    `json:"action"`
	Status    string    `json:"status"`
	Payload   T         `json:"payload"`
}

func NewGlobalEvent[T any](
	action, status, eventType string,
	payload T,
) GlobalEvent[T] {
	return GlobalEvent[T]{
		EventID:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now(),
		Source:    "payment-svc",
		Action:    action,
		Status:    status,
		Payload:   payload,
	}
}

func (ge GlobalEvent[T]) ToJSON() ([]byte, error) {
	return json.Marshal(ge)
}

func FromJSON[T any](data []byte) (GlobalEvent[T], error) {
	var ge GlobalEvent[T]
	err := json.Unmarshal(data, &ge)
	return ge, err
}
