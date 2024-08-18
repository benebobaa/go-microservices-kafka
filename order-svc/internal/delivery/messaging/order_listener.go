package messaging

import (
	"log"
	"order-svc/internal/dto"
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

		eventMsg, err := event.FromJSON[dto.OrderUpdateRequest, any](msg.Value)

		if err != nil {
			log.Println("failed parse event: ", err.Error())
		}

		switch eventMsg.State {
		case event.PAYMENT_SUCCESS.String():
			eventMsg.Payload.Request.Status = dto.COMPLETE.String()
			err = h.oc.UpdateOrderMessaging(sess.Context(), eventMsg)
			if err != nil {
				log.Println("failed update order: ", err.Error())
			}

		case event.PRODUCT_RESERVATION_FAILED.String(), event.PRODUCT_RELEASE_SUCCESS.String(), event.USER_VALIDATION_FAILED.String(), event.REFUND_SUCCESS.String():
			eventMsg.Payload.Request.Status = dto.CANCELLED.String()
			err = h.oc.UpdateOrderMessaging(sess.Context(), eventMsg)
			if err != nil {
				log.Println("failed update order: ", err.Error())
			}
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
