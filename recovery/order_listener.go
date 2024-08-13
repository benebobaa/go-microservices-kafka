package kafka

import (
	"fmt"
	"order-svc/internal/dto/event"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type MessageHandler struct {
	producer *KafkaProducer
}

func NewMessageHandler(producer *KafkaProducer) *MessageHandler {
	return &MessageHandler{producer: producer}
}

func (h MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		ev, err := event.FromJSON(msg.Value)
		if err != nil {
			fmt.Printf("Error decoding message: %v\n", err)
			continue
		}

		payload, err := ev.GetPayloadAsUserValidateRequest()
		if err != nil {
			fmt.Printf("Error getting payload: %v\n", err)
		}

		if ev.Source == "user_service" {
			fmt.Printf("Received message: %v\n", ev)
			return nil
		}

		ev.Payload = event.UserValidateRequest{
			Username: payload.Username,
		}

		eventBytes, err := ev.ToJSON()

		err = h.producer.SendMessage(uuid.New().String(), eventBytes)

		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
