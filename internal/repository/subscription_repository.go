package repository

import (
	"WeatherSubs/internal/models"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(subscription *models.Subscription) error
	GetByUserID(userID uint) ([]models.Subscription, error)
	GetAll() ([]models.Subscription, error)
	Delete(id uint) error
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(subscription *models.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *subscriptionRepository) GetByUserID(userID uint) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Where("user_id = ?", userID).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepository) GetAll() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := r.db.Find(&subscriptions).Error
	return subscriptions, err
}

func (r *subscriptionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Subscription{}, id).Error
}
