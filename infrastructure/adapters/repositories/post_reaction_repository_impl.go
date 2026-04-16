package repositories

import (
	"gorm.io/gorm"

	"gv/domain/models"
	"gv/domain/ports/repositories"
)

type PostReactionRepositoryImpl struct {
	db *gorm.DB
}

func NewPostReactionRepository(db *gorm.DB) repositories.PostReactionRepository {
	return &PostReactionRepositoryImpl{db: db}
}

func (r *PostReactionRepositoryImpl) Create(reaction *models.PostReaction) error {
	return r.db.Create(reaction).Error
}

func (r *PostReactionRepositoryImpl) Delete(postID, userID string) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&models.PostReaction{}).Error
}

func (r *PostReactionRepositoryImpl) Exists(postID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.PostReaction{}).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *PostReactionRepositoryImpl) CountByPost(postID string) (int, error) {
	var count int64
	err := r.db.Model(&models.PostReaction{}).
		Where("post_id = ?", postID).
		Count(&count).Error
	return int(count), err
}
