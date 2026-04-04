package models

import (
	"time"
)

type GameStatus string

const (
	StatusNowPlaying GameStatus = "NOW_PLAYING"
	StatusBacklog    GameStatus = "BACKLOG"
	StatusWishlist   GameStatus = "WISHLIST"
)

type Game struct {
	ID            string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID        string     `gorm:"not null;type:varchar(36);index" json:"user_id"`
	Name          string     `gorm:"not null;type:varchar(100)" json:"name"`
	Description   string     `gorm:"type:text" json:"description"`
	CoverImageURL string     `gorm:"type:varchar(500)" json:"cover_image_url"`
	Status        GameStatus `gorm:"type:enum('NOW_PLAYING','BACKLOG','WISHLIST');default:'WISHLIST'" json:"status"`
	Completed     bool       `gorm:"default:false" json:"completed"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Game) TableName() string {
	return "games"
}
