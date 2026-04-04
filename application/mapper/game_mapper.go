package mapper

import (
	"time"

	"gv/application/dto"
	"gv/domain/models"
)

func ToGameResponse(game *models.Game) dto.GameResponse {
	return dto.GameResponse{
		ID:            game.ID,
		UserID:        game.UserID,
		Name:          game.Name,
		Description:   game.Description,
		CoverImageURL: game.CoverImageURL,
		Status:        string(game.Status),
		Completed:     game.Completed,
		CreatedAt:     game.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     game.UpdatedAt.Format(time.RFC3339),
	}
}

func ToGameResponseList(games []models.Game) []dto.GameResponse {
	result := make([]dto.GameResponse, len(games))
	for i, game := range games {
		result[i] = ToGameResponse(&game)
	}
	return result
}
