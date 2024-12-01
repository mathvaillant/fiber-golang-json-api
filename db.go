package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("storage.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the db")
	}

	db.AutoMigrate(&User{}, &Book{})

	return db
}
