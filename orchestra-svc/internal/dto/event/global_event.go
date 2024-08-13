package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type State int

const (
	PENDING State = iota
	ORDER_CREATED
	PRODUCT_RESERVE_FAILED
	PAYMENT_FAILED
	PAYMENT_SUCCESS
)

func (s State) String() string {
	return [...]string{"PENDING", "ORDER_CREATED", "PRODUCT_RESERVE_FAILED", "PAYMENT_FAILED", "PAYMENT_SUCCESS"}[s]
}

type GlobalEvent[T any] struct {
	EventID    string    `json:"event_id"`
	InstanceID int32     `json:"instance_id"`
	EventType  string    `json:"event_type"`
	State      string    `json:"state"`
	Timestamp  time.Time `json:"timestamp"`
	Source     string    `json:"source"`
	Action     string    `json:"action"`
	Status     string    `json:"status"`
	Payload    T         `json:"payload"`
}

func NewGlobalEvent[T any](
	action, status, eventType string,
	payload T,
) GlobalEvent[T] {
	return GlobalEvent[T]{
		EventID:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now(),
		Source:    "orchestra-svc",
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
