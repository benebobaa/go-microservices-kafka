package pkg

import (
	"crypto/rsa"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

const (
	privKeyPath = "keys/private.pem"
	pubKeyPath  = "keys/public.pem"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

type Config struct {
	Port           string
	DBDriver       string
	DBSource       string
	KafkaBroker    string
	OrchestraTopic string
	OrderTopic     string
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
		GroupID:        os.Getenv("GROUP_ID"),
	}
}

func InitializeKeys() error {
	// Load private key
	privateKeyPEM, err := os.ReadFile(privKeyPath)
	if err != nil {
		return err
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return err
	}

	// Load public key
	publicKeyPEM, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return err
	}

	return nil
}
