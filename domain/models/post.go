package models

import (
	"time"
)

type PostType string

const (
	PostTypeTip        PostType = "TIP"
	PostTypeDiscussion PostType = "DISCUSSION"
	PostTypeReview     PostType = "REVIEW"
	PostTypeQuestion   PostType = "QUESTION"
	PostTypeNews       PostType = "NEWS"
)

type Post struct {
	ID             string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID         string    `gorm:"not null;type:varchar(36);index" json:"user_id"`
	GameName       string    `gorm:"not null;type:varchar(100);index" json:"game_name"`
	Title          string    `gorm:"not null;type:varchar(200)" json:"title"`
	Content        string    `gorm:"type:text" json:"content"`
	PostType       PostType  `gorm:"type:enum('TIP','DISCUSSION','REVIEW','QUESTION','NEWS');default:'DISCUSSION'" json:"post_type"`
	ReactionsCount int       `gorm:"default:0" json:"reactions_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Post) TableName() string {
	return "posts"
}
