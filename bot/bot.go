package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"mysub/internal/notify"
	"mysub/internal/storage"
	"mysub/internal/subscription"
	"strconv"
)

func InitBot(bot *tgbotapi.BotAPI) {
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	notify.Start(bot)
}

func ListenUpdates(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID

		if subscription.InProcess(chatID) {
			state := subscription.GetState(chatID)
			var response string
			switch state {
			case subscription.StepWaitingForTitle, subscription.StepWaitingForPrice, subscription.StepWaitingForDate:
				response, _ = subscription.ProcessSubscription(chatID, update.Message.Text)
			case subscription.StepWaitingForDelete:
				response = subscription.ProcessDelete(chatID, update.Message.Text)
			}
			bot.Send(tgbotapi.NewMessage(chatID, response))
			continue
		}

		if !update.Message.IsCommand() {
			bot.Send(tgbotapi.NewMessage(chatID, "Пожалуйста, используйте команду или завершите текущий ввод"))
			continue
		}

		switch update.Message.Command() {
		case "start":
			text := "Бот помогает отслеживать статус подписок на сервисы. \nИспользуй команды:" +
				"\n/setsub — указать подписку" +
				"\n/status — список подписок" +
				"\n/next — следующая дата списания" +
				"\n/delete — удалить подписку"
			bot.Send(tgbotapi.NewMessage(chatID, text))
		case "help":
			text := "Бот помогает отслеживать статус подписок на сервисы. \nИспользуй команды:" +
				"\n/setsub — указать подписку" +
				"\n/status — список подписок" +
				"\n/next — следующая дата списания" +
				"\n/delete — удалить подписку"
			bot.Send(tgbotapi.NewMessage(chatID, text))

		case "setsub":
			msg := subscription.StartSubscription(chatID)
			bot.Send(tgbotapi.NewMessage(chatID, msg))

		case "status":
			subs, err := storage.GetSubscriptionsByTelegramID(chatID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при получении подписок"))
				break
			}
			if len(subs) == 0 {
				bot.Send(tgbotapi.NewMessage(chatID, "У вас пока нет активных подписок."))
			} else {
				text := "Ваши подписки:\n"
				for i, sub := range subs {
					text += strconv.Itoa(i+1) + ". " + sub.Service + " — " +
						strconv.FormatFloat(sub.Price, 'f', 2, 64) + " руб. — следующее списание: " +
						sub.NextPayment.Format("2006-01-02") + "\n"
				}
				bot.Send(tgbotapi.NewMessage(chatID, text))
			}

		case "next":
			sub, err := storage.GetNextSubscription(chatID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при получении подписок"))
				break
			}

			if sub == nil {
				bot.Send(tgbotapi.NewMessage(chatID, "Не найдены активные подписки"))
			} else {
				text := "Ближайшее списание:\n"
				text += sub.Service + " — " +
					strconv.FormatFloat(sub.Price, 'f', 2, 64) + " руб. — следующее списание: " +
					sub.NextPayment.Format("2006-01-02") + "\n"
				bot.Send(tgbotapi.NewMessage(chatID, text))
			}

		case "delete":
			msg := subscription.StartDelete(chatID)
			bot.Send(tgbotapi.NewMessage(chatID, msg))

		default:
			bot.Send(tgbotapi.NewMessage(chatID, "Неизвестная команда"))
		}
	}
}
