package main

import (
	"log"
	"orchestra-svc/internal/app"
	"orchestra-svc/pkg"

	"github.com/gin-gonic/gin"
)

func main() {

	config := pkg.LoadConfig()

	log.Println("config: ", config)

	gin := gin.Default()
	db := pkg.NewDBConn(config.DBDriver, config.DBSource)

	app.NewApp(db, gin, config).Run()
}
