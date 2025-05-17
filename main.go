package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	tg "mysub/bot"
	"mysub/internal/storage"
	"os"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Printf("RAILWAY env: %v", os.Environ())
		log.Printf("DATABASE_URL: %s", dbURL)
		log.Fatal("DATABASE_URL не задан")
	}

	if err := storage.InitDb(dbURL); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatalf("Ошибка запуска бота: %v", err)
	}

	tg.InitBot(bot)
	tg.ListenUpdates(bot)

}
