package repositories

import (
	"gv/domain/models"
)

type PostReactionRepository interface {
	Create(reaction *models.PostReaction) error
	Delete(postID, userID string) error
	Exists(postID, userID string) (bool, error)
	CountByPost(postID string) (int, error)
}
