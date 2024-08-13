package messaging

import (
	"log"
	"orchestra-svc/internal/dto/event"
	"orchestra-svc/internal/usecase"

	"github.com/IBM/sarama"
)

type MessageHandler struct {
	oc *usecase.OrchestraUsecase
}

func NewMessageHandler(oc *usecase.OrchestraUsecase) *MessageHandler {
	return &MessageHandler{oc: oc}
}

func (h MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		eventMsg, err := event.FromJSON[any](msg.Value)

		if err != nil {
			log.Println("Error when parse message: ", err.Error())
		}

		err = h.oc.ProcessWorkflow(sess.Context(), eventMsg)

		if err != nil {
			log.Println("Error when process workflow: ", err.Error())
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
