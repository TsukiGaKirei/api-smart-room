package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type MySQLConnString struct {
	dbUsername             string
	dbPassword             string
	dbName                 string
	cloudSqlConnectionName string
}

func getMySQLCred() MySQLConnString {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error to load .env file!")
	}

	return MySQLConnString{
		dbUsername:             os.Getenv("DB_USERNAME"),
		dbPassword:             os.Getenv("DB_PASSWORD"),
		dbName:                 os.Getenv("DB_DATABASE"),
		cloudSqlConnectionName: os.Getenv("CLOUD_SQL_CONNECTION_NAME"),
	}
}

func GetMySQLConnString() string {
	credential := getMySQLCred()
	database := fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?parseTime=true",
		credential.dbUsername,
		credential.dbPassword,
		credential.cloudSqlConnectionName,
		credential.dbName,
	)

	return database
}
