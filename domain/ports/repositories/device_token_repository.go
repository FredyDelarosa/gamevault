package repositories

import (
	"gv/domain/models"
)

type DeviceTokenRepository interface {
	SaveToken(token *models.DeviceToken) error
	FindByUserID(userID string) ([]models.DeviceToken, error)
	DeleteByToken(fcmToken string) error
}