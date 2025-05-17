package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	tg "mysub/bot"
	"mysub/internal/storage"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TELEGRAM_APITOKEN")
	dsn := os.Getenv("DATABASE_URL")
	err = storage.InitDb(dsn)

	var bot *tgbotapi.BotAPI
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	tg.InitBot(bot)
	tg.ListenUpdates(bot)

}
