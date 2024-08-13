package main

import (
	"log"
	"order-svc/internal/app"
	"order-svc/pkg"

	"github.com/gin-gonic/gin"
)

func main() {

	config := pkg.LoadConfig()

	err := pkg.InitializeKeys()

	if err != nil {
		log.Fatal("Error initializing keys")
	}

	g := gin.Default()
	db := pkg.NewDBConn(config.DBDriver, config.DBSource)

	app.NewApp(db, g, config).Run()
}
