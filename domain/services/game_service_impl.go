package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"gv/core/logger"
	"gv/domain/models"
	"gv/domain/ports/repositories"
)

type GameServiceImpl struct {
	gameRepo repositories.GameRepository
}

func NewGameService(gameRepo repositories.GameRepository) *GameServiceImpl {
	return &GameServiceImpl{
		gameRepo: gameRepo,
	}
}

func (s *GameServiceImpl) CreateGame(userID string, name, description, coverImageURL string, status models.GameStatus) (*models.Game, error) {
	if name == "" {
		return nil, errors.New("game name is required")
	}

	game := &models.Game{
		ID:            uuid.New().String(),
		UserID:        userID,
		Name:          name,
		Description:   description,
		CoverImageURL: coverImageURL,
		Status:        status,
		Completed:     false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.gameRepo.Create(game); err != nil {
		logger.Error("Failed to create game: %v", err)
		return nil, errors.New("failed to create game")
	}

	return game, nil
}

func (s *GameServiceImpl) GetGamesByUser(userID string, status *models.GameStatus) ([]models.Game, error) {
	games, err := s.gameRepo.FindByUserID(userID, status)
	if err != nil {
		logger.Error("Failed to get games for user %s: %v", userID, err)
		return nil, errors.New("failed to retrieve games")
	}

	return games, nil
}

func (s *GameServiceImpl) GetGameByID(id, userID string) (*models.Game, error) {
	game, err := s.gameRepo.FindByID(id)
	if err != nil {
		logger.Error("Game not found: %s", id)
		return nil, errors.New("game not found")
	}

	// Validar que el juego exista
	if game == nil {
		logger.Error("Game not found: %s", id)
		return nil, errors.New("game not found")
	}

	if game.UserID != userID {
		logger.Error("User %s does not own game %s", userID, id)
		return nil, errors.New("access denied")
	}

	return game, nil
}

func (s *GameServiceImpl) UpdateGame(id, userID string, name, description, coverImageURL *string, status *models.GameStatus, completed *bool) (*models.Game, error) {
	game, err := s.gameRepo.FindByID(id)
	if err != nil {
		logger.Error("Game not found: %s", id)
		return nil, errors.New("game not found")
	}

	// Validar que el juego exista
	if game == nil {
		logger.Error("Game not found: %s", id)
		return nil, errors.New("game not found")
	}

	// Verificar propiedad
	if game.UserID != userID {
		logger.Error("User %s does not own game %s", userID, id)
		return nil, errors.New("access denied")
	}

	// Actualizar campos
	if name != nil {
		game.Name = *name
	}
	if description != nil {
		game.Description = *description
	}
	if coverImageURL != nil {
		game.CoverImageURL = *coverImageURL
	}
	if status != nil {
		game.Status = *status
	}
	if completed != nil {
		game.Completed = *completed
	}
	game.UpdatedAt = time.Now()

	if err := s.gameRepo.Update(game); err != nil {
		logger.Error("Failed to update game %s: %v", id, err)
		return nil, errors.New("failed to update game")
	}

	return game, nil
}

func (s *GameServiceImpl) DeleteGame(id, userID string) error {
	game, err := s.gameRepo.FindByID(id)
	if err != nil {
		logger.Error("Game not found: %s", id)
		return errors.New("game not found")
	}

	// Validar que el juego exista
	if game == nil {
		logger.Error("Game not found: %s", id)
		return errors.New("game not found")
	}
	
	if game.UserID != userID {
		logger.Error("User %s does not own game %s", userID, id)
		return errors.New("access denied")
	}

	if err := s.gameRepo.Delete(id); err != nil {
		logger.Error("Failed to delete game %s: %v", id, err)
		return errors.New("failed to delete game")
	}

	return nil
}
