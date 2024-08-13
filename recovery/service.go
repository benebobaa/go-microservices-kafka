package app

import (
	"order-svc/internal/delivery/http"
	kafka "order-svc/internal/delivery/messaging"
	"order-svc/internal/repository/sqlc"
)

func (app *App) startService() error {

	sqlc := sqlc.New(app.db)

	orchestraProducer, err := kafka.NewKafkaProducer(
		[]string{app.config.KafkaBroker},
		app.config.OrchestraTopic,
	)

	app.msg = kafka.NewMessageHandler(orchestraProducer)

	if err != nil {
		return err
	}

	authHanlder := http.NewAuthHandler()
	orderHandler := http.NewOrderHandler(sqlc, orchestraProducer)

	apiV1 := app.gin.Group("/api/v1")
	authV1 := apiV1.Group("/auth")
	orderV1 := apiV1.Group("/order")

	authHanlder.RegisterRoutes(authV1)
	orderHandler.RegisterRoutes(orderV1)

	return nil
}
