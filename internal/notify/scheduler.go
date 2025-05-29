package notify

import (
	"log"
	"mysub/internal/storage"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(bot *tgbotapi.BotAPI) {
	go func() {
		for {
			users, err := storage.GetAllTelegramIDs()
			if err != nil {
				log.Printf("ошибка получения списка пользователей: %v", err)
				time.Sleep(1 * time.Hour)
				continue
			}

			for _, chatID := range users {
				subs, err := storage.GetSubscriptionsByTelegramID(chatID)
				if err != nil {
					log.Printf("ошибка получения подписок пользователя %d: %v", chatID, err)
					continue
				}

				now := time.Now()
				for _, sub := range subs {
					if sub.NextPayment.Sub(now).Hours() <= 24 && sub.NextPayment.After(now) {
						msg := tgbotapi.NewMessage(chatID,
							"⏰ Напоминание: завтра будет списание по подписке «"+sub.Service+"» ("+
								sub.NextPayment.Format("02-01-2006")+")"+").")
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("ошибка отправки уведомления пользователю %d: %v", chatID, err)
						}
					}
				}
			}
			time.Sleep(12 * time.Hour)
		}
	}()
}
