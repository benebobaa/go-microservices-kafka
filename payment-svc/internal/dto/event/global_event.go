package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type State int

const (
	PENDING State = iota
	PRODUCT_RESERVATION_SUCCESS
	PRODUCT_RELEASE_SUCCESS
	ORDER_CANCEL
)

func (s State) String() string {
	return [...]string{"pending", "product_reservation_success", "product_release_success", "order_cancel"}[s]
}

type EventType int

const (
	ORDER_PROCESS EventType = iota
	ORDER_CANCEL_PROCESS
	BANK_ACCOUNT_REGISTRATION
)

func (e EventType) String() string {
	return [...]string{"order_process", "order_cancel_process", "bank_account_registration"}[e]
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
		Source:    "payment-svc",
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
