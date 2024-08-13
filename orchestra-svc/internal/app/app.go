package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"orchestra-svc/internal/delivery/messaging"
	"orchestra-svc/pkg"
	"orchestra-svc/pkg/consumer"
	"orchestra-svc/pkg/producer"
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
	msg    *messaging.MessageHandler
}

func NewApp(db *sql.DB, gin *gin.Engine, config *pkg.Config) *App {
	return &App{db: db, gin: gin, config: config}
}

func (app *App) Run() {

	//userProducer, err := producer.NewKafkaProducer(
	//	[]string{app.config.KafkaBroker},
	//	app.config.UserTopic,
	//)
	//if err != nil {
	//	log.Fatalf("Error creating Kafka producer: %v", err)
	//}
	//defer userProducer.Close()

	userProductProducer, err := producer.NewKafkaProducer(
		[]string{app.config.KafkaBroker},
		app.config.UserProductTopic,
	)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer userProductProducer.Close()

	//productProducer, err := producer.NewKafkaProducer(
	//	[]string{app.config.KafkaBroker},
	//	app.config.ProductTopic,
	//)
	//if err != nil {
	//	log.Fatalf("Error creating Kafka producer: %v", err)
	//}
	//defer productProducer.Close()

	if err := app.startService(userProductProducer); err != nil {
		log.Fatalf("Error starting service: %v", err)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", app.config.Port),
		Handler: app.gin,
	}

	consumer, err := consumer.NewKafkaConsumer(
		[]string{app.config.KafkaBroker},
		app.config.GroupID,
		[]string{app.config.OrchestraTopic},
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
		log.Printf("Starting server on port %s", app.config.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown Server ...")

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
