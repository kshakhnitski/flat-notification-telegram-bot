package main

import (
	"flat_bot/internal/bot"
	"flat_bot/internal/model"
	"flat_bot/internal/monitor"
	"flat_bot/internal/repository"
	"flat_bot/pkg/config"
	"flat_bot/pkg/db"
	"log"
	"time"
)

func main() {
	appConfig := config.LoadConfig()

	log.Println("Connecting to database...")
	dbConnection, err := db.ConnectDatabase(appConfig.DatabaseConfig)
	if err != nil {
		log.Fatalf("Error while connecting to database: %v", err)
	}
	log.Println("Connected to database")

	log.Println("Auto-migrating database schema...")
	err = dbConnection.AutoMigrate(&model.Flat{}, &model.User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}
	log.Println("Database schema migrated")

	flatRepository := repository.NewFlatRepository(dbConnection)
	userRepository := repository.NewUserRepository(dbConnection)

	telegramBot := bot.NewTelegramBot(&appConfig.TelegramBotConfig, userRepository)

	log.Println("Starting bot...")
	go telegramBot.Start()
	log.Println("Bot started")

	flatMonitor := monitor.NewFlatMonitor(30*time.Second, flatRepository, func(newFlats []model.Flat) {
		if len(newFlats) == 0 {
			log.Println("No new flats found")
			return
		}

		log.Printf("New flats found: %d. Notifying users.", len(newFlats))
		telegramBot.NotifyAboutNewFlats(newFlats)
	})

	log.Println("Starting flat monitor...")
	go flatMonitor.Start()
	log.Println("Flat monitor started")

	select {}
}
