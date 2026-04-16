package services

import (
	"gv/domain/models"
)

type PostService interface {
	CreatePost(userID, gameName, title, content string, postType models.PostType) (*models.Post, error)
	GetAllPosts(limit, offset int) ([]models.Post, error)
	GetPostsForUserGames(userID string, limit, offset int) ([]models.Post, error)
	GetPostByID(id string) (*models.Post, error)
	DeletePost(id, userID string) error
	ToggleReaction(postID, userID string) (bool, error)
	HasUserReacted(postID, userID string) (bool, error)
}
