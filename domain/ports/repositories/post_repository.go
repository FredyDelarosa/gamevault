package repositories

import (
	"gv/domain/models"
)

type PostRepository interface {
	Create(post *models.Post) error
	FindByID(id string) (*models.Post, error)
	FindAll(limit, offset int) ([]models.Post, error)
	FindByGameNames(gameNames []string, limit, offset int) ([]models.Post, error)
	FindByUserID(userID string) ([]models.Post, error)
	Delete(id string) error
	IncrementReactions(id string) error
	DecrementReactions(id string) error
}
