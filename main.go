package main // docker run --rm -v "$PWD":/app -w /app golang:1.25-alpine go run main.go -test

import (
	"flag"
	"time"

	"github.com/AgamariFF/TenderMessage.git/config"
	"github.com/AgamariFF/TenderMessage.git/internal/logger"
	"github.com/AgamariFF/TenderMessage.git/internal/scheduler"
	"github.com/AgamariFF/TenderMessage.git/internal/telegram"
	"github.com/AgamariFF/TenderMessage.git/internal/utils"
	"github.com/AgamariFF/TenderMessage.git/test"
)

func main() {
	testMode := flag.Bool("test", false, "Запуск в тестовом режиме")
	flag.Parse()
	if *testMode {
		test.Test()
		return
	}

	logger.InitLogger("debug")
	defer logger.Close()

	re, err := utils.LoadFilterPatterns("filter_patterns_vent.txt")
	if err != nil {
		logger.SugaredLogger.Errorf(err.Error())
	}

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.SugaredLogger.Errorf(err.Error())
	}

	teleBot, err := telegram.NewTelegramNotifier(cfg.TelegramToken, cfg.TelegramChatID)
	if err != nil {
		logger.SugaredLogger.Errorf(err.Error())
	}

	sched := scheduler.New(
		1*time.Hour,
		scheduler.Search,
		9,  // начало дня
		21, // конец дня
	)

	sched.Start(re, teleBot)
}
