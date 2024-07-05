package handler

import (
	"errors"
	"flat_bot/internal/model"
	"flat_bot/internal/repository"
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"os"
)

type StartHandler struct {
	Endpoint       string
	userRepository repository.UserRepository
}

func NewStartHandler(userRepository repository.UserRepository) *StartHandler {
	return &StartHandler{
		Endpoint:       "/start",
		userRepository: userRepository,
	}
}

func (h *StartHandler) Handle(c telebot.Context) error {
	sender := c.Sender()
	user := model.User{
		ID:        sender.ID,
		Username:  sender.Username,
		FirstName: sender.FirstName,
		ChatID:    c.Chat().ID,
	}

	log.Println("Start command from user: ", user)

	exists, err := h.userRepository.ExistsByID(user.ID)
	if err != nil {
		return err
	}

	if !exists {
		user, err = h.userRepository.Create(user)
		if err != nil {
			return err
		}
	}

	messagePattern, err := os.ReadFile("templates/hello.html")
	if err != nil {
		return errors.New(fmt.Sprintf("Error while reading hello message template: %v", err))
	}

	return c.Send(fmt.Sprintf(string(messagePattern), user.FirstName), telebot.ModeHTML)
}
