package messaging

import (
	"fmt"
	"payment-svc/internal/usecase"

	"github.com/IBM/sarama"
)

type MessageHandler struct {
	usecase *usecase.Usecase
}

func NewMessageHandler(usecase *usecase.Usecase) *MessageHandler {
	return &MessageHandler{usecase: usecase}
}

func (h MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Println("Received message", string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}
