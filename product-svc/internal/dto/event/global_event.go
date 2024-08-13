package event

import (
	"encoding/json"
	"time"
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

func (ge GlobalEvent[T]) ToJSON() ([]byte, error) {
	return json.Marshal(ge)
}

func (ge *GlobalEvent[T]) FromJSON(data []byte) error {
	return json.Unmarshal(data, ge)
}
