package mapper

import (
	"gv/application/dto"
	"gv/domain/models"
)

func ToUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}
