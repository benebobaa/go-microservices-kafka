package app

import (
	"orchestra-svc/internal/delivery/http"
	"orchestra-svc/internal/delivery/messaging"
	"orchestra-svc/internal/repository/sqlc"
	"orchestra-svc/internal/usecase"
	"orchestra-svc/pkg/producer"
)

func (app *App) startService(producer *producer.KafkaProducer) error {

	s := sqlc.New(app.db)
	// oc := usecase.NewOrderUsecase(producer)
	orc := usecase.NewOrchestraUsecase(s, producer)

	app.msg = messaging.NewMessageHandler(orc)

	wfu := usecase.NewWorkflowUsecase(s)
	wfh := http.NewWorkflowHandler(wfu)

	wfGroupV1 := app.gin.Group("/api/v1/workflow")
	wfh.RegisterRoutes(wfGroupV1)

	return nil
}
