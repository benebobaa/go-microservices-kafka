package event

import (
	"encoding/json"
	"fmt"
	"order-svc/pkg"
	"time"

	"github.com/google/uuid"
)

type State int

const (
	PENDING State = iota
	PRODUCT_RESERVE_FAILED
	PAYMENT_FAILED
	PAYMENT_SUCCESS
)

func (s State) String() string {
	return [...]string{"PENDING", "PRODUCT_RESERVE_FAILED", "PAYMENT_FAILED", "PAYMENT_SUCCESS"}[s]
}

type GlobalEvent[T any] struct {
	EventID    string    `json:"event_id"`
	InstanceID string    `json:"instance_id"`
	EventType  string    `json:"event_type"`
	State      string    `json:"state"`
	Timestamp  time.Time `json:"timestamp"`
	Source     string    `json:"source"`
	Action     string    `json:"action"`
	Status     string    `json:"status"`
	StatusCode int       `json:"status_code"`
	Payload    T         `json:"payload"`
}

func NewGlobalEvent[T any](
	action, status, eventType string,
	payload T,
) GlobalEvent[T] {
	return GlobalEvent[T]{
		EventID:    uuid.New().String(),
		InstanceID: fmt.Sprintf("I-%s", pkg.GenerateRandom6Char()),
		EventType:  eventType,
		Timestamp:  time.Now(),
		Source:     "order-svc",
		Action:     action,
		Status:     status,
		Payload:    payload,
	}
}

func FromJSON[T any](data []byte) (GlobalEvent[T], error) {
	var ge GlobalEvent[T]
	err := json.Unmarshal(data, &ge)
	return ge, err
}

func (ge GlobalEvent[T]) ToJSON() ([]byte, error) {
	return json.Marshal(ge)
}
