package models

import (
	"time"
)

type DeviceToken struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string    `gorm:"not null;type:varchar(36);index" json:"user_id"`
	FcmToken  string    `gorm:"not null;type:varchar(500);uniqueIndex" json:"fcm_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DeviceToken) TableName() string {
	return "device_tokens"
}