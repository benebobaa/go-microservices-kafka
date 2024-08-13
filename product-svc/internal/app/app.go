package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"product-svc/internal/delivery/kafka"
	"product-svc/internal/usecase"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	gin     *gin.Engine
	usecase *usecase.Usecase
}

func NewApp(gin *gin.Engine) *App {
	return &App{gin: gin}
}

func (a *App) Run() {
	brokers := []string{"localhost:29092"} // Replace with your Kafka broker addresses
	groupID := "user-svc-group"
	topics := []string{"product-topic"}
	orchestraTopic := "orchestra-topic"

	a.startService()

	server := http.Server{
		Addr:    ":8084",
		Handler: a.gin,
	}

	producer, err := kafka.NewKafkaProducer(brokers, orchestraTopic)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	consumer, err := kafka.NewKafkaConsumer(
		brokers, groupID, topics,
		a.usecase, producer,
	)
	defer consumer.Close()

	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	ctxCancel, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	go func() {
		if err := consumer.Consume(ctxCancel); err != nil {
			log.Fatalf("Error consuming Kafka messages: %v", err)
		}
	}()

	go func() {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Waiting signal send to chan quit
	// Blocking channel
	<-quit
	log.Println("Shutdown Server ...")
	log.Println("Closing Kafka consumer...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 1 seconds.")
	}

	log.Println("Server exiting")
}
