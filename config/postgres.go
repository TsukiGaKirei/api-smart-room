package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type PostgresCredential struct {
	DBUsername             string
	DBPassword             string
	DBName                 string
	DBHost                 string
	InstanceConnectionName string
}

func GetPostgresCredential() PostgresCredential {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error to load .env file")
	}

	return PostgresCredential{
		DBUsername:             os.Getenv("DB_USERNAME"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBName:                 os.Getenv("DB_DATABASE"),
		DBHost:                 os.Getenv("DB_HOST"),
		InstanceConnectionName: os.Getenv("CLOUD_SQL_CONNECTION_NAME"),
	}
}

func GetPostgresConnectionString() string {
	credential := GetPostgresCredential()
	// dataBase := fmt.Sprintf("user=%s password=%s database=%s host=/cloudsql/%s",
	dataBase := fmt.Sprintf("user=%s password=%s database=%s host=%s",

		credential.DBUsername,
		credential.DBPassword,
		credential.DBName,
		credential.DBHost,

		// credential.InstanceConnectionName,
	)
	return dataBase
}
