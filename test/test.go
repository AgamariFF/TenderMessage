package test

import (
	"fmt"
	"os"

	"github.com/AgamariFF/TenderMessage.git/config"
	"github.com/AgamariFF/TenderMessage.git/internal/logger"
	"github.com/AgamariFF/TenderMessage.git/internal/scheduler"
	"github.com/AgamariFF/TenderMessage.git/internal/telegram"
	"github.com/AgamariFF/TenderMessage.git/internal/utils"

	"github.com/joho/godotenv"
)

func Test() {
	logger.InitLogger("debug")
	defer logger.Close()
	logger.SugaredLogger.Infoln("=== ТЕСТ ПАРСЕРА И TELEGRAM ===")

	if err := godotenv.Load(); err != nil {
		logger.SugaredLogger.Infof("Не удалось загрузить .env файл: %v", err)
		logger.SugaredLogger.Infoln("Создайте файл .env с содержимым:")
		logger.SugaredLogger.Infoln("TELEGRAM_TOKEN=ваш_токен")
		logger.SugaredLogger.Infoln("TELEGRAM_CHAT_ID=ваш_chat_id")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.SugaredLogger.Errorln(err.Error())
	}

	if cfg.TelegramToken == "" || cfg.TelegramChatID == 0 {
		logger.SugaredLogger.Errorln("Токен или Chat ID не указаны в .env файле")
	}

	re, err := utils.LoadFilterPatterns("filter_patterns_vent.txt")
	if err != nil {
		logger.SugaredLogger.Errorln(err.Error())
	}

	logger.SugaredLogger.Infoln("Парсим тендеры с сайта...")

	tenders, err := scheduler.Search(re)

	if err != nil {
		logger.SugaredLogger.Errorf("Ошибка парсинга: %v", err)
	}

	logger.SugaredLogger.Infof("Найдено тендеров: %d\n", len(tenders))

	if len(tenders) == 0 {
		logger.SugaredLogger.Infoln("Тендеры не найдены. Попробуйте изменить параметры поиска.")
		return
	}

	logger.SugaredLogger.Infoln("\n📋 НАЙДЕННЫЕ ТЕНДЕРЫ:")
	for i, tender := range tenders {
		fmt.Printf("\n%d. %s\n", i+1, tender.Title)
		fmt.Printf("   Заказчик: %s\n", tender.Customer)
		fmt.Printf("   Цена: %s\n", tender.Price)
		fmt.Printf("   Дата публикации: %s\n", tender.PublishDate)
		fmt.Printf("   Регион: %s\n", tender.Region)
		fmt.Printf("   Ссылка: %s\n", tender.Link)
	}

	logger.SugaredLogger.Infoln("\nОтправляем в Telegram...")

	notifier, err := telegram.NewTelegramNotifier(cfg.TelegramToken, int64(cfg.TelegramChatID))
	if err != nil {
		logger.SugaredLogger.Errorf("Ошибка создания телеграм-бота: %v", err)
	}

	// Отправляем только первые 3 тендера для теста
	maxTenders := 3
	if len(tenders) < maxTenders {
		maxTenders = len(tenders)
	}

	err = notifier.SendTenderNotification(tenders[:maxTenders])
	if err != nil {
		logger.SugaredLogger.Errorf("Ошибка отправки в Telegram: %v", err)
	}

	logger.SugaredLogger.Infof("Тест завершен! Отправлено %d тендеров в Telegram.\n", maxTenders)
	logger.SugaredLogger.Infoln("Проверьте Telegram-бота.")
}
