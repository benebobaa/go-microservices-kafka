package app

import (
	"product-svc/internal/delivery/messaging"
	"product-svc/internal/provider"
	"product-svc/internal/usecase"
	"product-svc/pkg/http_client"
	"product-svc/pkg/producer"
	"time"
)

func (app *App) startService(orchestraProducer *producer.KafkaProducer) error {

	productClient := http_client.NewProductClient(
		app.config.ClientUrl,
		5*time.Second,
	)

	productProvider := provider.NewProductProviderImpl(productClient)

	u := usecase.NewUsecase(productProvider, orchestraProducer)

	app.msg = messaging.NewMessageHandler(u)

	return nil
}
