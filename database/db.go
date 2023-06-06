package database

import (
	"api-smart-room/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func Init(params ...string) *gorm.DB {
	conString := config.GetPostgresConnectionString()
	log.Print(conString)

	var err error
	db, err = gorm.Open(postgres.Open(conString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Panic(err)
	}

	return db
}

func GetDBInstance() *gorm.DB {
	return db
}
