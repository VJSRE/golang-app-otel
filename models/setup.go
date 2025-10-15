package models

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {

	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {

		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = database.AutoMigrate(&Book{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = database
}
