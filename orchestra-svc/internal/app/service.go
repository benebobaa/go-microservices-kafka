package app

import (
	"orchestra-svc/internal/delivery/http"
	"orchestra-svc/internal/delivery/messaging"
	"orchestra-svc/internal/repository/cache"
	"orchestra-svc/internal/repository/sqlc"
	"orchestra-svc/internal/usecase"
	"orchestra-svc/pkg/producer"
)

func (app *App) startService(userProductProducer *producer.KafkaProducer) error {

	s := sqlc.New(app.db)
	c := cache.NewPayloadCache()
	// oc := usecase.NewOrderUsecase(producer)
	orc := usecase.NewOrchestraUsecase(s, userProductProducer, c)

	app.msg = messaging.NewMessageHandler(orc)

	wfh := http.NewWorkflowHandler()

	wfGroupV1 := app.gin.Group("/api/v1/workflow")
	wfh.RegisterRoutes(wfGroupV1)

	return nil
}
