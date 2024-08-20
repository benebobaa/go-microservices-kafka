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

	s := sqlc.NewStore(app.db)
	c := cache.NewPayloadCache()

	orc := usecase.NewOrchestraUsecase(s, userProductProducer, c)
	rc := usecase.NewRetryUsecase(s, userProductProducer, orc)

	app.msg = messaging.NewMessageHandler(orc)

	wfh := http.NewWorkflowHandler(orc, rc)

	wfGroupV1 := app.gin.Group("/api/v1/workflow")
	wfh.RegisterRoutes(wfGroupV1)

	return nil
}
