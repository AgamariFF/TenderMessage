package telegram

import (
	"fmt"
	"strings"
	"time"

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

	batchSize := 5
	for i := 0; i < len(tenders); i += batchSize {
		end := i + batchSize
		if end > len(tenders) {
			end = len(tenders)
		}

		batch := tenders[i:end]
		err := n.sendBatch(batch, i/batchSize+1, (len(tenders)+batchSize-1)/batchSize)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *TelegramNotifier) sendBatch(tenders []models.Tender, batchNum, totalBatches int) error {
	var builder strings.Builder

	if totalBatches > 1 {
		builder.WriteString(fmt.Sprintf("🚨 <b>Новые тендеры (часть %d/%d):</b>\n\n", batchNum, totalBatches))
	} else {
		builder.WriteString("🚨 <b>Новые тендеры:</b>\n\n")
	}

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

	_, err := n.bot.Send(
		&telebot.Chat{ID: n.chatID},
		builder.String(),
		&telebot.SendOptions{
			ParseMode:             telebot.ModeHTML,
			DisableWebPagePreview: true,
		},
	)

	return err
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
