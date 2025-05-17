package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"mysub/api"
	tg "mysub/bot"
	"mysub/internal/storage"
	"mysub/models"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = storage.InitDb()
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	var bot *tgbotapi.BotAPI
	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	tg.InitBot(bot)
	tg.ListenUpdates(bot)

	subs := []models.Subscription{
		{
			Service:  "Spotify",
			Price:    1.0,
			CreateDt: time.Now(),
			Category: "music",
		},
		{
			Service:  "AWS",
			Price:    1.0,
			CreateDt: time.Now(),
			Category: "cloud",
		},
	}

	r.Use(api.SetUserCtx(subs))

	r.Get("/api/user/me", api.HandleGet)

	http.ListenAndServe(":8080", r)
}
