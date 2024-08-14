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
	return [...]string{"PENDING", "PRODUCT_RESERVE_FAILED", "payment_failed", "payment_success"}[s]
}

type EventType int

const (
	ORDER_PROCESS State = iota
	ORDER_CANCEL_PROCESS
)

func (e EventType) String() string {
	return [...]string{"order_process", "order_cancel_process"}[e]
}

type BasePayload[R any, S any] struct {
	Request  R `json:"request"`
	Response S `json:"response"`
}

type GlobalEvent[R any, S any] struct {
	EventID    string            `json:"event_id"`
	InstanceID string            `json:"instance_id"`
	EventType  string            `json:"event_type"`
	State      string            `json:"state"`
	Timestamp  time.Time         `json:"timestamp"`
	Source     string            `json:"source"`
	Action     string            `json:"action"`
	Status     string            `json:"status"`
	StatusCode int               `json:"status_code"`
	Payload    BasePayload[R, S] `json:"payload"`
}

func NewGlobalEvent[R any, S any](
	action, status, eventType string,
	payload BasePayload[R, S],
) GlobalEvent[R, S] {
	return GlobalEvent[R, S]{
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
func FromJSON[R any, S any](data []byte) (GlobalEvent[R, S], error) {
	var ge GlobalEvent[R, S]
	err := json.Unmarshal(data, &ge)
	return ge, err
}

func (ge GlobalEvent[R, S]) ToJSON() ([]byte, error) {
	return json.Marshal(ge)
}
