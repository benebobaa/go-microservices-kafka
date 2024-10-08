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

		log.Println("Event type: ", eventMsg.EventType)
		log.Println("Event state: ", eventMsg.State)

		switch eventMsg.EventType {
		case event.ORDER_PROCESS.String():
			if eventMsg.State == event.PRODUCT_RESERVATION_SUCCESS.String() {
				err = h.u.ProcessPaymentMessaging(sess.Context(), eventMsg)
			}
		case event.ORDER_CANCEL_PROCESS.String():
			if eventMsg.State == event.PRODUCT_RELEASE_SUCCESS.String() {
				err = h.u.RefundPaymentMessaging(sess.Context(), eventMsg)
			}
		case event.BANK_ACCOUNT_REGISTRATION.String():
			eventMsg, _ := event.FromJSON[dto.AccountBalanceRequest, any](msg.Value)
			err = h.u.CreateAccountBalanceMessaging(sess.Context(), eventMsg)
		}

		if err != nil {
			log.Println("Error processing message: ", err)
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
