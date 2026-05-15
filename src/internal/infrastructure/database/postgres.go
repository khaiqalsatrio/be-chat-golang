package database

import (
	"chat-golang/src/config"
	"chat-golang/src/internal/domain/entities"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db

	// Auto Migration
	err = DB.AutoMigrate(
		&entities.User{},
		&entities.Message{},
		&entities.Room{},
		&entities.Agenda{},
		&entities.Status{},
		&entities.Post{},
		&entities.Like{},
		&entities.Comment{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	log.Println("Database connection established and migrated")
}
