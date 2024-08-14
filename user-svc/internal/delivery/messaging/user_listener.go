package messaging

import (
	"log"
	"user-svc/internal/dto"
	"user-svc/internal/dto/event"
	"user-svc/internal/usecase"

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

		eventMsg, err := event.FromJSON[dto.UserValidateRequest, any](msg.Value)
		if err != nil {
			log.Println("Error when parse message: ", err.Error())
		}

		err = h.u.ValidateUserMessaging(sess.Context(), eventMsg)
		if err != nil {
			log.Println("Error when validate user: ", err.Error())
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
