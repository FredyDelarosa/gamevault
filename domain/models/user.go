package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null;type:varchar(100)" json:"email"`
	Password  string    `gorm:"not null;type:varchar(255)" json:"-"`
	FirstName string    `gorm:"not null;type:varchar(50)" json:"first_name"`
	LastName  string    `gorm:"not null;type:varchar(50)" json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Games     []Game    `gorm:"foreignKey:UserID" json:"games,omitempty"`
}

func (User) TableName() string {
	return "users"
}
