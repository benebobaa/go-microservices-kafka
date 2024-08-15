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
	PAYMENT_SUCCESS State = iota
	PRODUCT_RELEASE_SUCCESS
	USER_VALIDATION_FAILED
)

func (s State) String() string {
	return [...]string{"payment_success", "product_release_success", "user_validation_failed"}[s]
}

type EventType int

const (
	ORDER_PROCESS EventType = iota
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
	action, status, state, eventType string,
	payload BasePayload[R, S],
) GlobalEvent[R, S] {
	return GlobalEvent[R, S]{
		EventID:    uuid.New().String(),
		InstanceID: fmt.Sprintf("I-%s", pkg.GenerateRandom6Char()),
		State:      state,
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
