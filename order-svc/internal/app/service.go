package app

import (
	"order-svc/internal/delivery/http"
	"order-svc/internal/delivery/messaging"
	"order-svc/internal/middleware"
	"order-svc/internal/repository/sqlc"
	"order-svc/internal/usecase"
	"order-svc/pkg/producer"
)

func (app *App) startService(orchestraProducer *producer.KafkaProducer) error {

	sqlc := sqlc.NewStore(app.db)

	orderUsecase := usecase.NewOrderUsecase(sqlc, orchestraProducer)
	bankRegisUsecase := usecase.NewBankRegistrationUsecase(sqlc, orchestraProducer)

	app.msg = messaging.NewMessageHandler(orderUsecase, bankRegisUsecase)

	authHandler := http.NewAuthHandler()
	orderHandler := http.NewOrderHandler(orderUsecase)
	bankRegisHandler := http.NewBankRegisHandler(bankRegisUsecase)

	apiV1 := app.gin.Group("/api/v1")

	authV1 := apiV1.Group("/auth")
	orderV1 := apiV1.Group("/order")
	bankRegisV1 := apiV1.Group("/bank")

	orderV1.Use(middleware.AuthMiddleware())

	authHandler.RegisterRoutes(authV1)
	orderHandler.RegisterRoutes(orderV1)
	bankRegisHandler.RegisterRoutes(bankRegisV1)

	return nil
}
