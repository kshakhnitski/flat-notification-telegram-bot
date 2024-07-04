package main

import (
	"flat_bot/internal/bot"
	"flat_bot/internal/listener"
	"flat_bot/internal/model"
	"flat_bot/internal/repository"
	"flat_bot/pkg/config"
	"flat_bot/pkg/db"
	"log"
	"time"
)

func main() {
	appConfig := config.LoadConfig()

	dbConnection, err := db.ConnectDatabase(appConfig.DatabaseConfig)
	if err != nil {
		log.Fatalf("Error while connecting to database: %v", err)
	}

	err = dbConnection.AutoMigrate(&model.Flat{}, &model.User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	flatRepository := repository.NewFlatRepository(dbConnection)
	userRepository := repository.NewUserRepository(dbConnection)

	telegramBot := bot.NewTelegramBot(&appConfig.TelegramBotConfig, userRepository)

	log.Println("Starting bot...")
	go telegramBot.Start()
	log.Println("Bot started")

	flatListener := listener.NewFlatListener(30*time.Second, flatRepository, func(newFlats []model.Flat) {
		log.Printf("New flats found: %d. Notifying users.", len(newFlats))
		for _, flat := range newFlats {
			telegramBot.NotifyAboutNewFlat(flat)
		}
	})

	log.Println("Starting flat listener...")
	go flatListener.Start()
	log.Println("Flat listener started")

	select {}
}
