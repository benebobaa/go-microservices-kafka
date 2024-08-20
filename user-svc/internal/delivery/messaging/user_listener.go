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

		eventMsg, err := event.FromJSON[any, any](msg.Value)
		if err != nil {
			log.Println("Error when parse message: ", err.Error())
		}

		log.Println("Event type: ", eventMsg.EventType)
		log.Println("Event state: ", eventMsg.State)

		switch eventMsg.EventType {
		case event.BANK_ACCOUNT_REGISTRATION.String():
			if eventMsg.State == event.BANK_BALANCE_CREATED.String() {
				eventMsg, _ := event.FromJSON[dto.UpdateBankIDRequest, any](msg.Value)
				err = h.u.UpdateUserMessaging(sess.Context(), eventMsg)
			} else {
				eventMsg, _ := event.FromJSON[dto.UserCreateRequest, any](msg.Value)
				err = h.u.CreateUserMessaging(sess.Context(), eventMsg)
			}

		case event.ORDER_PROCESS.String():
			eventMsg, _ := event.FromJSON[dto.UserValidateRequest, any](msg.Value)
			err = h.u.UserDetailMessaging(sess.Context(), eventMsg)

		case event.ORDER_CANCEL_PROCESS.String():
			eventMsg, _ := event.FromJSON[dto.UserValidateRequest, any](msg.Value)
			err = h.u.UserDetailMessaging(sess.Context(), eventMsg)
		}

		if err != nil {
			log.Println("Error when validate user: ", err.Error())
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
