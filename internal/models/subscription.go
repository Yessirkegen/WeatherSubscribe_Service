package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	City      string         `json:"city"`
	Timezone  string         `json:"timezone"`
	Type      string         `json:"type"`      // Тип уведомления: email, sms и т.д.
	Frequency string         `json:"frequency"` // Частота: daily, weekly и т.д.
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
