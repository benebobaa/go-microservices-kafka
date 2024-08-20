package app

import (
	"time"
	"user-svc/internal/delivery/messaging"
	"user-svc/internal/provider"
	"user-svc/internal/usecase"
	"user-svc/pkg/http_client"
	"user-svc/pkg/producer"
)

func (app *App) startService(orchestraProducer *producer.KafkaProducer) error {

	userClient := http_client.NewUserClient(
		app.config.ClientUrl,
		5*time.Second,
	)

	userProvider := provider.NewUserProvider(userClient)

	uc := usecase.NewUsecase(userClient, orchestraProducer, userProvider)

	app.msg = messaging.NewMessageHandler(uc)

	return nil
}
