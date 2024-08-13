package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"payment-svc/internal/app"
	"payment-svc/pkg"
)

func main() {

	config := pkg.LoadConfig()

	log.Println("config: ", config)

	g := gin.New()
	app.NewApp(g, config).Run()
}
