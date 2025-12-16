package scheduler

import (
	"log"
	"regexp"
	"time"

	"github.com/AgamariFF/TenderMessage.git/internal/logger"
	"github.com/AgamariFF/TenderMessage.git/internal/models"
	"github.com/AgamariFF/TenderMessage.git/internal/telegram"
)

type SimpleScheduler struct {
	interval     time.Duration
	task         func(*regexp.Regexp) ([]models.Tender, error)
	lastResults  []models.Tender // Результаты предыдущего выполнения
	dayStartHour int
	dayEndHour   int
}

func New(interval time.Duration, task func(*regexp.Regexp) ([]models.Tender, error), startHour, endHour int) *SimpleScheduler {
	return &SimpleScheduler{
		interval:     interval,
		task:         task,
		lastResults:  []models.Tender{},
		dayStartHour: startHour,
		dayEndHour:   endHour,
	}
}

func (s *SimpleScheduler) isDayTime() bool {
	now := time.Now()
	currentHour := now.Hour()

	return currentHour >= s.dayStartHour && currentHour < s.dayEndHour
}

// getNewTenders находит тендеры, которых не было в предыдущих результатах
func (s *SimpleScheduler) getNewTenders(currentTenders []models.Tender) []models.Tender {
	newTenders := []models.Tender{}

	oldTendersMap := make(map[string]bool)
	for _, tender := range s.lastResults {
		oldTendersMap[tender.Link] = true
	}

	for _, tender := range currentTenders {
		if !oldTendersMap[tender.Link] {
			newTenders = append(newTenders, tender)
		}
	}

	return newTenders
}

// Start запускает планировщик (блокирующая функция)
func (s *SimpleScheduler) Start(re *regexp.Regexp, teleBot *telegram.TelegramNotifier) {
	log.Println("Планировщик запущен. Интервал:", s.interval)

	tenders, err := s.task(re)
	if err != nil {
		logger.SugaredLogger.Warnf("Ошибка выполнения task: %s\n", err)
	}

	s.lastResults = tenders

	for {
		time.Sleep(s.interval)

		if !s.isDayTime() {
			log.Printf("Ночное время (%02d). Пропускаем выполнение.\n", time.Now().Hour())
			continue
		}

		log.Println("Выполняем проверку...")

		currentTenders, err := s.task(re)
		if err != nil {
			logger.SugaredLogger.Warnf("Ошибка выполнения task: %s\n", err)
		}

		newTenders := s.getNewTenders(currentTenders)

		if len(newTenders) > 0 {
			logger.SugaredLogger.Infof("Найдено новых тендеров: %d\n", len(newTenders))

			err := teleBot.SendTenderNotification(newTenders)
			if err != nil {
				log.Printf("Ошибка отправки в Telegram: %v", err)
			}

			s.lastResults = currentTenders

			for _, tender := range newTenders {
				logger.SugaredLogger.Infof("НОВЫЙ ТЕНДЕР: %s - %s\n", tender.Link, tender.Title)
			}
		} else {
			logger.SugaredLogger.Debugln("Новых тендеров не найдено")
			// Все равно обновляем lastResults, на случай если какие-то тендеры были удалены
			s.lastResults = currentTenders
		}
	}
}
