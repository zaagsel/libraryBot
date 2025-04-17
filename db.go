package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var errDB error
	db, errDB = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDB != nil {
		log.Fatal("Error connecting to database:", errDB)
	}

	log.Println("Successfully connected to database!")
}

func MigrateDB() {
	// Автоматическое создание таблиц
	err := db.AutoMigrate(&User{}, &Genre{}, &Status{}, &Book{})
	if err != nil {
		log.Fatal("Error migrating database:", err)
	}

	log.Println("Database tables created/updated successfully!")
}
