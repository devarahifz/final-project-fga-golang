package database

import (
	"final-project/helpers"
	"final-project/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	var (
		host     = helpers.GetEnv("DB_HOST")
		port     = helpers.GetEnv("DB_PORT")
		user     = helpers.GetEnv("DB_USERNAME")
		password = helpers.GetEnv("DB_PASSWORD")
		dbname   = helpers.GetEnv("DB_NAME")
	)

	config := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	db.AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})
}

func GetDB() *gorm.DB {
	return db
}
