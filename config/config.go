package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken  string
	TelegramChatID int64
}

func LoadConfig(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, err
	}

	chatID, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	return &Config{
		TelegramToken:  os.Getenv("TELEGRAM_TOKEN"),
		TelegramChatID: chatID,
	}, nil
}
