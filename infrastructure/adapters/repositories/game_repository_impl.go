package repositories

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"gv/domain/models"
	"gv/domain/ports/repositories"
)

type GameRepositoryImpl struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) repositories.GameRepository {
	return &GameRepositoryImpl{db: db}
}

func (r *GameRepositoryImpl) Create(game *models.Game) error {
	return r.db.Create(game).Error
}

func (r *GameRepositoryImpl) FindByID(id string) (*models.Game, error) {
	var game models.Game
	err := r.db.Where("id = ?", id).First(&game).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &game, err
}

func (r *GameRepositoryImpl) FindByUserID(userID string, status *models.GameStatus) ([]models.Game, error) {
	var games []models.Game
	query := r.db.Where("user_id = ?", userID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("created_at DESC").Find(&games).Error
	return games, err
}

func (r *GameRepositoryImpl) Update(game *models.Game) error {
	return r.db.Save(game).Error
}

func (r *GameRepositoryImpl) Delete(id string) error {
	return r.db.Delete(&models.Game{}, "id = ?", id).Error
}

func (r *GameRepositoryImpl) FindByNameMatch(gameName string) ([]models.Game, error) {
	var games []models.Game
	err := r.db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(gameName)+"%").Find(&games).Error
	return games, err
}
