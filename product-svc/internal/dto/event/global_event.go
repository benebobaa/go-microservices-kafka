package event

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type State int

const (
	USER_VALIDATION_SUCCESS State = iota
	PRODUCT_RELEASE_SUCCESS
	PAYMENT_FAILED
	ORDER_CANCEL
	PRODUCT_RETRY
	REFUND_FAILED
)

func (s State) String() string {
	return [...]string{"user_validation_success", "product_release_success", "payment_failed", "order_cancel", "product_retry", "refund_failed"}[s]
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
	action, status, state string,
	payload BasePayload[R, S],
) GlobalEvent[R, S] {
	return GlobalEvent[R, S]{
		EventID:   uuid.New().String(),
		State:     state,
		Timestamp: time.Now(),
		Source:    "product-svc",
		Action:    action,
		Status:    status,
		Payload:   payload,
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
