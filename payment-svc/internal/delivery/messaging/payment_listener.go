package messaging

import (
	"log"
	"payment-svc/internal/dto"
	"payment-svc/internal/dto/event"
	"payment-svc/internal/usecase"

	"github.com/IBM/sarama"
)

type MessageHandler struct {
	u *usecase.Usecase
}

func NewMessageHandler(u *usecase.Usecase) *MessageHandler {
	return &MessageHandler{u: u}
}

func (h MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		eventMsg, err := event.FromJSON[dto.PaymentRequest, any](msg.Value)

		if err != nil {
			log.Println("Error parsing message: ", err)
		}

		err = h.u.ProcessPaymentMessaging(sess.Context(), eventMsg)
		if err != nil {
			log.Println("Error processing message: ", err)
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
