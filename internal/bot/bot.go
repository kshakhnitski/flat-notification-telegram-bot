package bot

import (
	"flat_bot/internal/bot/handler"
	"flat_bot/internal/model"
	"flat_bot/internal/repository"
	"flat_bot/pkg/config"
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"os"
	"time"
)

type TelegramBot struct {
	bot            *telebot.Bot
	userRepository repository.UserRepository
}

func NewTelegramBot(config *config.TelegramBotConfig, userRepository repository.UserRepository) *TelegramBot {
	pref := telebot.Settings{
		Token:  config.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatalf("Error while creating bot: %v", err)
	}

	initializeHandlers(bot, userRepository)

	return &TelegramBot{bot: bot, userRepository: userRepository}
}

func initializeHandlers(bot *telebot.Bot, userRepository repository.UserRepository) {
	startHandler := handler.NewStartHandler(userRepository)
	bot.Handle(startHandler.Endpoint, startHandler.Handle)
}

func (b *TelegramBot) Start() {
	b.bot.Start()
}

func (b *TelegramBot) NotifyAboutNewFlats(flats []model.Flat) {
	messagePattern, err := os.ReadFile("templates/new_flat_available.html")
	if err != nil {
		log.Printf("Error while reading message template: %v", err)
		return
	}

	users, err := b.userRepository.FindAll()
	if err != nil {
		log.Printf("Error while getting users from database: %v", err)
		return
	}

	for _, flat := range flats {
		if flat.Metro == nil {
			flat.Metro = new(string)
			*flat.Metro = "Не указано"
		}

		message := fmt.Sprintf(
			string(messagePattern),
			flat.ID,
			flat.Source,
			flat.Parameters,
			flat.Address,
			flat.Description,
			*flat.Metro,
			flat.Link,
			flat.PriceInUsd,
			flat.PriceInByn,
		)

		for _, user := range users {
			_, err := b.bot.Send(telebot.ChatID(user.ChatID), message, telebot.ModeHTML)
			if err != nil {
				log.Printf("Error while sending message to user: %v", err)
			}
		}
	}

}
