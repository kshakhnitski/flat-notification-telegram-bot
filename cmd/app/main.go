package main

import (
	"flat_bot/internal/bot"
	"flat_bot/internal/model"
	"flat_bot/internal/parser"
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

	flatRepository := repository.NewFlatRepository(dbConnection)
	userRepository := repository.NewUserRepository(dbConnection)

	telegramBot := bot.NewTelegramBot(&appConfig.TelegramBotConfig, userRepository)

	log.Println("Starting bot...")
	go telegramBot.Start()
	log.Println("Bot started")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			log.Println("Checking for new flats...")
			newFlats := checkForNewFlats(flatRepository)

			if len(newFlats) == 0 {
				continue
			}

			log.Printf("New flats found: %d. Notifying users.", len(newFlats))
			for _, flat := range newFlats {
				telegramBot.NotifyAboutNewFlat(flat)
			}
		}
	}()

	select {}
}

func checkForNewFlats(flatRepository repository.FlatRepository) []model.Flat {
	var loadedFlats []model.Flat
	var newFlats []model.Flat

	loadedFlats = append(loadKufarFlats())

	for _, flat := range loadedFlats {
		exists, err := flatRepository.ExistsByID(flat.ID)
		if err != nil {
			log.Fatalf("Error while checking if flat exists: %v", err)
		}

		if exists {
			continue
		}

		newFlats = append(newFlats, flat)
		_, err = flatRepository.Create(flat)
		if err != nil {
			log.Fatalf("Error while saving flat: %v", err)
		}
	}

	return newFlats
}

func loadKufarFlats() []model.Flat {
	url := "https://re.kufar.by/l/minsk/snyat/kvartiru-dolgosrochno/bez-posrednikov?cur=USD&gbx=b%3A27.150276810254205%2C53.344712700318226%2C28.932808548535448%2C54.058539029097716&prc=r%3A0%2C350&rms=v.or%3A1%2C2%2C3&size=30"
	kufarFlatParser := parser.NewKufarFlatParser(url)

	flats, err := kufarFlatParser.Parse()
	if err != nil {
		log.Fatalf("Error while parsing kufar flats: %v", err)
	}

	return flats
}
