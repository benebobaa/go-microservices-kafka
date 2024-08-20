package app

import (
	"payment-svc/internal/delivery/messaging"
	"payment-svc/internal/provider"
	"payment-svc/internal/usecase"
	"payment-svc/pkg/http_client"
	"payment-svc/pkg/producer"
	"time"
)

func (app *App) startService(orchestraProducer *producer.KafkaProducer) error {

	client := http_client.NewPaymentClient(
		app.config.ClientUrl,
		5*time.Second,
	)

	paymentProvider := provider.NewPaymentProviderImpl(client)

	uc := usecase.NewUsecase(paymentProvider, orchestraProducer)

	app.msg = messaging.NewMessageHandler(uc)

	return nil
}
