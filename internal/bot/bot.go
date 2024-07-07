package bot

import (
	"flat_bot/internal/bot/handler"
	"flat_bot/internal/model"
	"flat_bot/internal/repository"
	"flat_bot/pkg/config"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"time"
)

type TelegramBot struct {
	bot            *tele.Bot
	userRepository repository.UserRepository
}

func NewTelegramBot(config *config.TelegramBotConfig, userRepository repository.UserRepository) (*TelegramBot, error) {
	pref := tele.Settings{
		Token:  config.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	b := &TelegramBot{bot: bot, userRepository: userRepository}
	b.initHandlers()

	return b, nil
}

func (b TelegramBot) initHandlers() {
	startHandler := handler.NewStartHandler(b.userRepository)
	b.bot.Handle(startHandler.Endpoint, startHandler.Handle)
}

func (b TelegramBot) Start() {
	b.bot.Start()
}

func (b TelegramBot) NotifyAboutNewFlats(flats []model.Flat) {
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
			_, err := b.bot.Send(tele.ChatID(user.ChatID), message, tele.ModeHTML)
			if err != nil {
				log.Printf("Error while sending message to user: %v", err)
			}
		}
	}

}
