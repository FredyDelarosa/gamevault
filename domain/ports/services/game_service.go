package services

import (
	"gv/domain/models"
)

type GameService interface {
	CreateGame(userID string, name, description, coverImageURL string, status models.GameStatus) (*models.Game, error)
	GetGamesByUser(userID string, status *models.GameStatus) ([]models.Game, error)
	GetGameByID(id, userID string) (*models.Game, error)
	UpdateGame(id, userID string, name, description, coverImageURL *string, status *models.GameStatus, completed *bool) (*models.Game, error)
	DeleteGame(id, userID string) error
}
