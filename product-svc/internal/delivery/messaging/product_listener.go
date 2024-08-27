package messaging

import (
	"github.com/IBM/sarama"
	"log"
	"product-svc/internal/dto"
	"product-svc/internal/dto/event"
	"product-svc/internal/usecase"
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

		eventMsg, err := event.FromJSON[dto.ProductRequest, any](msg.Value)

		if err != nil {
			log.Println("Error parsing message: ", err)
		}

		switch eventMsg.EventType {
		case event.ORDER_PROCESS.String():
			if eventMsg.State == event.USER_VALIDATION_SUCCESS.String() {
				err = h.u.ReserveProductMessaging(sess.Context(), eventMsg)
			}

			if eventMsg.State == event.PAYMENT_FAILED.String() {
				err = h.u.ReleaseProductMessaging(sess.Context(), eventMsg)
			}

			if eventMsg.State == event.PRODUCT_RETRY.String() {
				err = h.u.ReserveProductMessaging(sess.Context(), eventMsg)
			}

			if err != nil {
				log.Println("Error processing message: ", err)
			}

		case event.ORDER_CANCEL_PROCESS.String():
			if eventMsg.State == event.USER_VALIDATION_SUCCESS.String() {
				err = h.u.ReleaseProductMessaging(sess.Context(), eventMsg)
			}

			if eventMsg.State == event.REFUND_FAILED.String() {
				err = h.u.ReserveProductMessaging(sess.Context(), eventMsg)
			}

			if err != nil {
				log.Println("Error processing message: ", err)
			}
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
