package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBDriver       string
	DBSource       string
	KafkaBroker    string
	OrchestraTopic string
	OrderTopic     string
	UserTopic      string
	ProductTopic   string
	GroupID        string
}

func LoadConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Port:           os.Getenv("PORT"),
		DBDriver:       os.Getenv("DB_DRIVER"),
		DBSource:       os.Getenv("DB_SOURCE"),
		KafkaBroker:    os.Getenv("KAFKA_BROKER"),
		OrchestraTopic: os.Getenv("ORCHESTRA_TOPIC"),
		OrderTopic:     os.Getenv("ORDER_TOPIC"),
		UserTopic:      os.Getenv("USER_TOPIC"),
		ProductTopic:   os.Getenv("PRODUCT_TOPIC"),
		GroupID:        os.Getenv("GROUP_ID"),
	}
}
