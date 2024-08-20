package main

import (
	"fmt"
	"mock-svc/handler"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello, World!")

	gin := gin.Default()

	uh := handler.NewUserHandler()
	ph := handler.NewProductHandler()
	pyh := handler.NewPaymentHandler()

	gin.GET("/users", uh.GetUser)
	gin.GET("/users/:username", uh.GetuUserByUsername)
	gin.POST("/users", uh.CreateUser)
	gin.PATCH("/users/:username", uh.UpdateUserBankID)

	gin.GET("/products", ph.GetProduct)
	gin.POST("/products/reserve", ph.ReserveProduct)
	gin.POST("/products/release", ph.ReleaseProduct)

	gin.GET("/balances", pyh.GetBalance)
	gin.GET("/transactions", pyh.GetTransaction)
	gin.POST("/payment", pyh.CreateTransaction)
	gin.PATCH("/payment/refund", pyh.RefundTransaction)
	gin.POST("/payment/balances", pyh.CreateBalance)

	gracefulShutdown(pyh)

	fmt.Println("Server is running on port 5000")
	gin.Run(":5000")
}

func gracefulShutdown(ph *handler.PaymentHandler) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nShutting down server...")

		if err := ph.SaveData(); err != nil {
			fmt.Println("Error saving data:", err)
		}

		os.Exit(0)
	}()
}
