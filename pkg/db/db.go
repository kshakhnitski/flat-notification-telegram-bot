package db

import (
	"flat_bot/internal/model"
	"flat_bot/pkg/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectDatabase(config config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Europe/Minsk",
		config.Host, config.Port, config.User, config.Password, config.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Flat{}, &model.User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	return db, nil
}
