package repositories

import (
	"gv/domain/models"
)

type GameRepository interface {
	Create(game *models.Game) error
	FindByID(id string) (*models.Game, error)
	FindByUserID(userID string, status *models.GameStatus) ([]models.Game, error)
	Update(game *models.Game) error
	Delete(id string) error
	FindByNameMatch(gameName string) ([]models.Game, error) // <-- NUEVO
}
