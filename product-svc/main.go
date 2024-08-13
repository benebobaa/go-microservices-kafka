package main

import (
	"product-svc/internal/app"

	"github.com/gin-gonic/gin"
)

func main() {
	gin := gin.New()
	app.NewApp(gin).Run()
}
