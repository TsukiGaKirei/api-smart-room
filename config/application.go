package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetApplicationPort() string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error to load .env file")
	}

	return os.Getenv("APP_PORT")
}

func GetSignatureKey() []byte {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error to load .env file")
	}

	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

func GetAppUrl() string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error to load .env file")
	}

	return os.Getenv("APP_URL") + ":" + os.Getenv("APP_PORT")
}
