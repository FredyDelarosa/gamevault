package services

import (
	"gv/domain/models"
)

type AuthService interface {
	Register(email, password, firstName, lastName string) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
	ValidateToken(tokenString string) (string, error)
}
