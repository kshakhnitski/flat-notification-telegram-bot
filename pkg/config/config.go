package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	TelegramBotConfig TelegramBotConfig
	DatabaseConfig    DatabaseConfig
}

type TelegramBotConfig struct {
	Token string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func LoadConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Error parsing DB_PORT: %v", err)
	}

	return AppConfig{
		TelegramBotConfig: TelegramBotConfig{
			Token: os.Getenv("TELEGRAM_BOT_API_KEY"),
		},
		DatabaseConfig: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
		},
	}
}
