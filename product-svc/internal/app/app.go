package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"product-svc/internal/delivery/messaging"
	"product-svc/pkg"
	"product-svc/pkg/consumer"
	"product-svc/pkg/producer"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	gin    *gin.Engine
	config *pkg.Config
	msg    *messaging.MessageHandler
}

func NewApp(gin *gin.Engine, c *pkg.Config) *App {
	return &App{
		gin:    gin,
		config: c,
	}
}

func (app *App) Run() {

	orchestraProducer, err := producer.NewKafkaProducer(
		[]string{app.config.KafkaBroker},
		app.config.OrchestraTopic,
	)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer orchestraProducer.Close()

	app.startService(orchestraProducer)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", app.config.Port),
		Handler: app.gin,
	}

	consumer, err := consumer.NewKafkaConsumer(
		[]string{app.config.KafkaBroker},
		app.config.GroupID,
		[]string{app.config.ProductTopic, app.config.UserProductTopic},
		app.msg,
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
