package messaging

import (
	"log"
	"order-svc/internal/dto"
	"order-svc/internal/dto/event"
	"order-svc/internal/usecase"

	"github.com/IBM/sarama"
)

type MessageHandler struct {
	oc  *usecase.OrderUsecase
	brc *usecase.BankRegistrationUsecase
}

func NewMessageHandler(oc *usecase.OrderUsecase, brc *usecase.BankRegistrationUsecase) *MessageHandler {
	return &MessageHandler{
		oc:  oc,
		brc: brc,
	}
}

func (h MessageHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		eventMsg, err := event.FromJSON[dto.OrderUpdateRequest, any](msg.Value)

		if err != nil {
			log.Println("failed parse event: ", err.Error())
		}

		switch eventMsg.EventType {
		case event.ORDER_PROCESS.String():
			if eventMsg.State == event.PAYMENT_SUCCESS.String() {
				eventMsg.Payload.Request.Status = dto.COMPLETE.String()
				err = h.oc.UpdateOrderMessaging(sess.Context(), eventMsg)
			}

			if eventMsg.State == event.USER_VALIDATION_FAILED.String() {
				eventMsg.Payload.Request.Status = dto.CANCELLED.String()
				err = h.oc.UpdateOrderMessaging(sess.Context(), eventMsg)
			}

			if eventMsg.State == event.PRODUCT_RESERVATION_FAILED.String() {
				eventMsg.Payload.Request.Status = dto.CANCELLED.String()
				err = h.oc.UpdateOrderMessaging(sess.Context(), eventMsg)
			}

			if eventMsg.State == event.PRODUCT_RELEASE_SUCCESS.String() {
				eventMsg.Payload.Request.Status = dto.CANCELLED.String()
				err = h.oc.UpdateOrderMessaging(sess.Context(), eventMsg)
			}

		case event.ORDER_CANCEL_PROCESS.String():

			if eventMsg.State == event.USER_VALIDATION_FAILED.String() {
				eventMsg.Payload.Request.Status = dto.COMPLETE.String()
				err = h.oc.UpdateOrderMessaging(sess.Context(), eventMsg)
			}

			if eventMsg.State == event.REFUND_FAILED.String() {
				eventMsg.Payload.Request.Status = dto.COMPLETE.String()
				err = h.oc.UpdateOrderMessaging(sess.Context(), eventMsg)
			}

			if eventMsg.State == event.REFUND_SUCCESS.String() {
				eventMsg.Payload.Request.Status = dto.CANCELLED.String()
				err = h.oc.UpdateOrderMessaging(sess.Context(), eventMsg)
			}

		case event.BANK_ACCOUNT_REGISTRATION.String():
			if eventMsg.State == event.USER_BANKID_UPDATED.String() {
				eventMsg, _ := event.FromJSON[dto.BankRegistrationUpdate, any](msg.Value)
				eventMsg.Payload.Request.Status = dto.COMPLETE.String()
				err = h.brc.UpdateBankRegistrationMessaging(sess.Context(), eventMsg)
			}
		}

		if err != nil {
			log.Println("failed process event: ", err.Error())
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
