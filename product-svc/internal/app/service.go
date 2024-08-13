package app

import (
	"product-svc/http_client"
	"product-svc/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *App) startService() error {

	userClient := http_client.NewUserClient(
		"https://beneboba.wiremockapi.cloud/users/validate",
		5*time.Second,
	)

	usecase := usecase.NewUsecase(userClient)
	app.usecase = usecase

	app.gin.GET("/products", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	return nil
}
