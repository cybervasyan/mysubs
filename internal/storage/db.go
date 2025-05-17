package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"mysub/models"
)

var DB *gorm.DB

func InitDb(dbToken string) error {
	if dbToken == "" {
		return fmt.Errorf("DATABASE_URL не задан в окружении")
	}

	db, err := gorm.Open(postgres.Open(dbToken), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Не удалось подключиться к БД: %w", err)
	}

	DB = db
	_ = db.AutoMigrate(&models.Subscription{})
	return nil
}
