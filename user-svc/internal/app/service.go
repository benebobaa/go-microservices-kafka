package app

import (
	"time"
	"user-svc/internal/delivery/http"
	"user-svc/internal/delivery/messaging"
	"user-svc/internal/usecase"
	"user-svc/pkg/http_client"
	"user-svc/pkg/producer"
)

func (app *App) startService(orchestraProducer *producer.KafkaProducer) error {

	userClient := http_client.NewUserClient(
		app.config.ClientUrl,
		5*time.Second,
	)

	usecase := usecase.NewUsecase(userClient, orchestraProducer)

	app.msg = messaging.NewMessageHandler(usecase)

	handler := http.NewHandler(usecase)
	app.gin.GET("/send-message", handler.TestValidate)

	return nil
}
