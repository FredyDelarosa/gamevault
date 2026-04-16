package models

import (
	"time"
)

type PostReaction struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	PostID    string    `gorm:"not null;type:varchar(36);index;uniqueIndex:idx_post_user" json:"post_id"`
	UserID    string    `gorm:"not null;type:varchar(36);index;uniqueIndex:idx_post_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (PostReaction) TableName() string {
	return "post_reactions"
}
