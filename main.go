package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tg "mysub/bot"
	"mysub/internal/storage"
	"os"
)

func main() {
	token := os.Getenv("TELEGRAM_APITOKEN")
	dsn := os.Getenv("DATABASE_URL")
	err := storage.InitDb(dsn)

	var bot *tgbotapi.BotAPI
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	tg.InitBot(bot)
	tg.ListenUpdates(bot)

}
