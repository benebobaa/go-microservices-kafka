package messaging

import (
	"log"
	"order-svc/internal/dto/event"
	"order-svc/internal/usecase"

	"github.com/IBM/sarama"
)

type MessageHandler struct {
	oc *usecase.OrderUsecase
}

func NewMessageHandler(oc *usecase.OrderUsecase) *MessageHandler {
	return &MessageHandler{oc: oc}
}

func (h MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		eventMsg, err := event.FromJSON[any, any](msg.Value)

		if err != nil {
			log.Println("failed parse event: ", err.Error())
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
