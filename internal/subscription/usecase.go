package subscription

import (
	"log"
	"mysub/internal/storage"
	"mysub/models"
	"strconv"
	"time"
)

type Step string

const (
	StepWaitingForTitle  Step = "waiting_for_title"
	StepWaitingForPrice  Step = "waiting_for_price"
	StepWaitingForDate   Step = "waiting_for_date"
	StepWaitingForDelete Step = "waiting_for_delete"
)

var userState = make(map[int64]Step)
var tempSubs = make(map[int64]models.Subscription)

func StartSubscription(chatID int64) string {
	userState[chatID] = StepWaitingForTitle
	return "Введите название подписки:"
}

func ProcessSubscription(chatID int64, input string) (string, bool) {
	state, ok := userState[chatID]
	if !ok {
		return "", false
	}

	sub := tempSubs[chatID]
	switch state {
	case StepWaitingForTitle:
		sub.Service = input
		tempSubs[chatID] = sub
		userState[chatID] = StepWaitingForPrice
		return "Введите цену (в формате хх.хх):", true

	case StepWaitingForPrice:
		price, err := parsePrice(input)
		if err != nil {
			return "Неверный формат цены. Введите число (в формате xx.xx):", true
		}
		sub.Price = price
		tempSubs[chatID] = sub
		userState[chatID] = StepWaitingForDate
		return "Введите дату следующего списания (в формате DD.MM.YYYY):", true

	case StepWaitingForDate:
		date, err := time.Parse("02.01.2006", input)
		if err != nil {
			log.Printf("Неверный формат даты: %v", err)
			return "Неверный формат даты. Попробуйте ещё раз:", true
		}
		sub.NextPayment = date
		sub.CreateDt = time.Now()
		sub.TelegramID = chatID

		err = storage.SaveSubscription(&sub)
		if err != nil {
			log.Printf("ошибка сохранения в БД: %v", err)
			return "Ошибка при сохранении подписки", true
		}

		delete(userState, chatID)
		delete(tempSubs, chatID)
		return "Подписка успешно добавлена!", false
	}

	return "", false
}

func StartDelete(chatID int64) string {
	userState[chatID] = StepWaitingForDelete
	return "Введите название подписки, которую хотите удалить:"
}

func ProcessDelete(chatID int64, name string) string {
	sub, err := storage.DeleteSubscriptionByName(chatID, name)
	if err != nil {
		delete(userState, chatID)
		return "Подписка с таким названием не найдена или уже удалена."
	}

	subs, err := storage.GetSubscriptionsByTelegramID(chatID)
	if err != nil {
		delete(userState, chatID)
		return "Удаление прошло, но не удалось получить список подписок."
	}

	msg := "Удалена подписка: " + sub.Service + "\n\nТекущие подписки:\n"
	if len(subs) == 0 {
		msg += "Нет активных подписок."
	} else {
		for i, s := range subs {
			msg += strconv.Itoa(i+1) + ". " + s.Service + " — " +
				strconv.FormatFloat(s.Price, 'f', 2, 64) + " руб. — следующее списание: " +
				s.NextPayment.Format("2006-01-02") + "\n"
		}
	}
	delete(userState, chatID)
	return msg
}

func InProcess(chatID int64) bool {
	_, ok := userState[chatID]
	return ok
}

func parsePrice(input string) (float64, error) {
	return strconv.ParseFloat(input, 64)
}

func GetState(chatID int64) Step {
	return userState[chatID]
}
