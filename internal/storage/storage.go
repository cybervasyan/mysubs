package storage

import (
	"mysub/models"
)

func SaveSubscription(sub *models.Subscription) error {
	return DB.Create(sub).Error
}

func GetSubscriptionsByTelegramID(telegramID int64) ([]models.Subscription, error) {
	var subs []models.Subscription
	err := DB.Where("telegram_id = ?", telegramID).Find(&subs).Error
	return subs, err
}

func DeleteSubscriptionByName(telegramID int64, name string) (*models.Subscription, error) {
	var sub models.Subscription
	tx := DB.
		Where("telegram_id = ? AND service = ?", telegramID, name).
		First(&sub)
	if tx.Error != nil {
		return nil, tx.Error
	}

	err := DB.Delete(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func GetNextSubscription(telegramID int64) (*models.Subscription, error) {
	var sub models.Subscription
	tx := DB.
		Where("telegram_id = ?", telegramID).
		Order("next_payment desc").
		First(&sub)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &sub, nil
}
