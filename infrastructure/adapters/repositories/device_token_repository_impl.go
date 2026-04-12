package repositories

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gv/domain/models"
	"gv/domain/ports/repositories"
)

type DeviceTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewDeviceTokenRepository(db *gorm.DB) repositories.DeviceTokenRepository {
	return &DeviceTokenRepositoryImpl{db: db}
}

func (r *DeviceTokenRepositoryImpl) SaveToken(token *models.DeviceToken) error {
	// Upsert: si el fcm_token ya existe, actualizar el user_id y updated_at
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "fcm_token"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "updated_at"}),
	}).Create(token).Error
}

func (r *DeviceTokenRepositoryImpl) FindByUserID(userID string) ([]models.DeviceToken, error) {
	var tokens []models.DeviceToken
	err := r.db.Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

func (r *DeviceTokenRepositoryImpl) DeleteByToken(fcmToken string) error {
	return r.db.Where("fcm_token = ?", fcmToken).Delete(&models.DeviceToken{}).Error
}