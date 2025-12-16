package telegram

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AgamariFF/TenderMessage.git/internal/excel"
	"github.com/AgamariFF/TenderMessage.git/internal/logger"
	"github.com/AgamariFF/TenderMessage.git/internal/models"
	"gopkg.in/telebot.v3"
)

type TelegramNotifier struct {
	bot    *telebot.Bot
	chatID int64
}

func NewTelegramNotifier(token string, chatID int64) (*TelegramNotifier, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка создания бота: %w", err)
	}

	return &TelegramNotifier{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (n *TelegramNotifier) SendTenderNotification(tenders []models.Tender) error {
	if len(tenders) == 0 {
		return nil
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("🚨 <b>Новые тендеры :</b>"))

	for i, tender := range tenders {
		if i > 0 {
			builder.WriteString("\n────────────\n\n")
		}

		builder.WriteString(fmt.Sprintf(
			"<b>%d. %s</b>\n\n"+
				"👤 <b>Заказчик:</b> %s\n"+
				"💰 <b>Цена:</b> %s\n"+
				"📍 <b>Регион:</b> %s\n"+
				"🔗 <b>Ссылка:</b> <a href=\"%s\">открыть</a>\n",
			i+1,
			limitString(tender.Title, 200),
			limitString(tender.Customer, 120),
			tender.Price,
			limitString(tender.Region, 111),
			tender.Link,
		))
	}

	msg, err := n.bot.Send(
		&telebot.Chat{ID: n.chatID},
		builder.String(),
		&telebot.SendOptions{
			ParseMode:             telebot.ModeHTML,
			DisableWebPagePreview: true,
		},
	)

	filename, err := excel.ToExcel(&tenders)
	if err != nil {
		logger.SugaredLogger.Warnf(err.Error())
	}

	defer os.Remove(filename)

	fileCaption := fmt.Sprintf(
		"📊 <b>Полный отчет по %d тендерам</b>\n"+
			"📅 Дата: %s\n"+
			"💾 Файл: Excel (.xlsx)\n\n"+
			"Содержит полную информацию из уведомления выше.",
		len(tenders),
		time.Now().Format("02.01.2006 15:04"),
	)

	doc := &telebot.Document{
		File:     telebot.FromDisk(filename),
		FileName: fmt.Sprintf("tenders_%s.xlsx", time.Now().Format("20060102")),
		Caption:  fileCaption,
	}

	_, err = n.bot.Reply(msg, doc, &telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
	})

	if err != nil {
		logger.SugaredLogger.Warnf("Ошибка отправки Excel: %v", err)
		_, err = n.bot.Send(&telebot.Chat{ID: n.chatID}, doc, &telebot.SendOptions{
			ParseMode: telebot.ModeHTML,
		})
		return err
	}

	return nil
}

// Вспомогательная функция для ограничения длины строки
func limitString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func (n *TelegramNotifier) Start() {
	logger.SugaredLogger.Infoln("Telegram бот запущен...")
	n.bot.Start()
}
