package messaging

import (
	"context"
	"log"
	"orchestra-svc/internal/dto/event"
	"orchestra-svc/internal/usecase"

	"github.com/IBM/sarama"
)

type MessageHandler struct {
	usecase *usecase.OrchestraUsecase
}

func NewMessageHandler(u *usecase.OrchestraUsecase) *MessageHandler {
	return &MessageHandler{usecase: u}
}

func (h MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		eventMsg, err := event.FromJSON[any](msg.Value)

		if err != nil {
			log.Println("Error when parse message: ", err.Error())
		}

		if eventMsg.InstanceID != 0 {
			err = h.usecase.ProcessWorkflow(context.Background(), eventMsg)
		} else {
			err = h.usecase.ProcessNewWorkflow(context.Background(), eventMsg)
		}

		if err != nil {
			log.Println("Error when process workflow: ", err.Error())
			continue
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
