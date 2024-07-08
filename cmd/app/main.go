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
	cfg := config.LoadConfig()

	log.Println("Connecting to database...")
	conn, err := db.ConnectDatabase(cfg.DatabaseConfig)
	if err != nil {
		log.Fatalf("Error while connecting to database: %v", err)
	}
	log.Println("Connected to database")

	log.Println("Auto-migrating database schema...")
	err = conn.AutoMigrate(&model.Flat{}, &model.User{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}
	log.Println("Database schema migrated")

	flatRepo := repository.NewFlatRepository(conn)
	userRepo := repository.NewUserRepository(conn)

	flatBot, err := bot.NewTelegramBot(cfg.TelegramBotConfig, userRepo)
	if err != nil {
		log.Fatalf("Error while creating telegram bot: %v", err)
	}

	log.Println("Starting bot...")
	go flatBot.Start()
	log.Println("Bot started")

	flatMonitor := monitor.NewFlatMonitor(30*time.Second, flatRepo, func(newFlats []model.Flat) {
		if len(newFlats) == 0 {
			log.Println("No new flats found")
			return
		}

		log.Printf("New flats found: %d. Notifying users.", len(newFlats))
		flatBot.NotifyAboutNewFlats(newFlats)
	})

	log.Println("Starting flat monitor...")
	go flatMonitor.Start()
	log.Println("Flat monitor started")

	select {}
}
