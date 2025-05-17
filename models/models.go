package models

import "time"

type Subscription struct {
	ID          int64     `gorm:"primary_key" json:"id"`
	TelegramID  int64     `gorm:"index" json:"telegramId"`
	Service     string    `json:"service"`
	Price       float64   `json:"price"`
	CreateDt    time.Time `json:"createDt"`
	Category    string    `json:"category"`
	NextPayment time.Time `json:"nextPayment"`
}
