package app

import (
	"payment-svc/internal/delivery/http"
	"payment-svc/internal/delivery/messaging"
	"payment-svc/internal/usecase"
	"payment-svc/pkg/http_client"
	"payment-svc/pkg/producer"
	"time"
)

func (app *App) startService(orchestraProducer *producer.KafkaProducer) error {

	userClient := http_client.NewPaymentClient(
		app.config.ClientUrl,
		5*time.Second,
	)

	usecase := usecase.NewUsecase(userClient, orchestraProducer)

	app.msg = messaging.NewMessageHandler(usecase)

	_ = http.NewHandler(usecase)

	return nil
}
