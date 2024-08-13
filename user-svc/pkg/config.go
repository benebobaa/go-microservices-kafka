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
	UserTopic      string
	GroupID        string
	ClientUrl      string
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
		UserTopic:      os.Getenv("USER_TOPIC"),
		GroupID:        os.Getenv("GROUP_ID"),
		ClientUrl:      os.Getenv("CLIENT_URL"),
	}
}
