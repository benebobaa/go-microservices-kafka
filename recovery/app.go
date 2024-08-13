package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	kafka "order-svc/internal/delivery/messaging"
	"order-svc/pkg"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	db     *sql.DB
	gin    *gin.Engine
	config *pkg.Config
	msg    *kafka.MessageHandler
}

func NewApp(db *sql.DB, gin *gin.Engine, config *pkg.Config) *App {
	return &App{db: db, gin: gin, config: config}
}

func (a *App) Run() {

	a.startService()

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.PORT),
		Handler: a.gin,
	}

	consumer, err := kafka.NewKafkaConsumer(
		[]string{a.config.KafkaBroker},
		a.config.GroupID,
		[]string{a.config.OrderTopic},
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

	// Waiting signal send to chan quit
	// Blocking channel
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}

	log.Println("Server exiting")
}
