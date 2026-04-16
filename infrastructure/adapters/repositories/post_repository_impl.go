package repositories

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"gv/domain/models"
	"gv/domain/ports/repositories"
)

type PostRepositoryImpl struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) repositories.PostRepository {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepositoryImpl) FindByID(id string) (*models.Post, error) {
	var post models.Post
	err := r.db.Where("id = ?", id).First(&post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &post, err
}

func (r *PostRepositoryImpl) FindAll(limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error
	return posts, err
}

func (r *PostRepositoryImpl) FindByGameNames(gameNames []string, limit, offset int) ([]models.Post, error) {
	var posts []models.Post

	if len(gameNames) == 0 {
		return posts, nil
	}

	// Matching case-insensitive usando LOWER
	query := r.db
	conditions := []string{}
	args := []interface{}{}

	for _, name := range gameNames {
		conditions = append(conditions, "LOWER(game_name) LIKE ?")
		args = append(args, "%"+strings.ToLower(name)+"%")
	}

	whereClause := strings.Join(conditions, " OR ")
	err := query.Where(whereClause, args...).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&posts).Error

	return posts, err
}

func (r *PostRepositoryImpl) FindByUserID(userID string) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&posts).Error
	return posts, err
}

func (r *PostRepositoryImpl) Delete(id string) error {
	return r.db.Delete(&models.Post{}, "id = ?", id).Error
}

func (r *PostRepositoryImpl) IncrementReactions(id string) error {
	return r.db.Model(&models.Post{}).Where("id = ?", id).
		UpdateColumn("reactions_count", gorm.Expr("reactions_count + 1")).Error
}

func (r *PostRepositoryImpl) DecrementReactions(id string) error {
	return r.db.Model(&models.Post{}).Where("id = ? AND reactions_count > 0", id).
		UpdateColumn("reactions_count", gorm.Expr("reactions_count - 1")).Error
}
