package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	kafka "user-svc/internal/delivery/messaging"
	"user-svc/internal/usecase"
	"user-svc/pkg"
	"user-svc/pkg/consumer"
	"user-svc/pkg/producer"

	"github.com/gin-gonic/gin"
)

type App struct {
	gin     *gin.Engine
	usecase *usecase.Usecase
	config  *pkg.Config
	msg     *kafka.MessageHandler
}

func NewApp(gin *gin.Engine, c *pkg.Config) *App {
	return &App{
		gin:     gin,
		usecase: &usecase.Usecase{},
		config:  c,
	}
}

func (a *App) Run() {

	orchestraProducer, err := producer.NewKafkaProducer(
		[]string{a.config.KafkaBroker},
		a.config.OrchestraTopic,
	)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer orchestraProducer.Close()

	a.startService(orchestraProducer)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.Port),
		Handler: a.gin,
	}

	consumer, err := consumer.NewKafkaConsumer(
		[]string{a.config.KafkaBroker},
		a.config.GroupID,
		[]string{a.config.UserTopic},
		a.msg,
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
