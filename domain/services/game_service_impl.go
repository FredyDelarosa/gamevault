package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"gv/core/logger"
	"gv/domain/models"
	"gv/domain/ports/repositories"
	"gv/domain/ports/services"
)

type GameServiceImpl struct {
	gameRepo            repositories.GameRepository
	notificationService services.NotificationService
}

func NewGameService(gameRepo repositories.GameRepository, notificationService services.NotificationService) *GameServiceImpl {
	return &GameServiceImpl{
		gameRepo:            gameRepo,
		notificationService: notificationService,
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

	// Notificación automática al crear juego
	go s.notificationService.SendNotificationToUser(
		userID,
		"Juego agregado",
		"\""+name+"\" se añadió a tu biblioteca",
		"game_updates",
	)

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

	if game == nil {
		logger.Error("Game not found: %s", id)
		return nil, errors.New("game not found")
	}

	if game.UserID != userID {
		logger.Error("User %s does not own game %s", userID, id)
		return nil, errors.New("access denied")
	}

	// Guardar nombre antes de actualizar para la notificación
	gameName := game.Name

	if name != nil {
		game.Name = *name
		gameName = *name
	}
	if description != nil {
		game.Description = *description
	}
	if coverImageURL != nil {
		game.CoverImageURL = *coverImageURL
	}

	// Detectar cambio de estado para notificación específica
	var statusChanged bool
	var newStatus models.GameStatus
	if status != nil && *status != game.Status {
		statusChanged = true
		newStatus = *status
		game.Status = *status
	}

	// Detectar si se completó
	var justCompleted bool
	if completed != nil && *completed && !game.Completed {
		justCompleted = true
		game.Completed = *completed
	} else if completed != nil {
		game.Completed = *completed
	}

	game.UpdatedAt = time.Now()

	if err := s.gameRepo.Update(game); err != nil {
		logger.Error("Failed to update game %s: %v", id, err)
		return nil, errors.New("failed to update game")
	}

	// Notificaciones automáticas según el tipo de cambio
	if justCompleted {
		go s.notificationService.SendNotificationToUser(
			userID,
			"¡Juego completado!",
			"Felicidades, completaste \""+gameName+"\"",
			"game_updates",
		)
	} else if statusChanged {
		statusLabel := ""
		switch newStatus {
		case models.StatusNowPlaying:
			statusLabel = "Now Playing"
		case models.StatusBacklog:
			statusLabel = "Backlog"
		case models.StatusWishlist:
			statusLabel = "Wishlist"
		}
		go s.notificationService.SendNotificationToUser(
			userID,
			"Juego movido",
			"\""+gameName+"\" se movió a "+statusLabel,
			"game_updates",
		)
	}

	return game, nil
}

func (s *GameServiceImpl) DeleteGame(id, userID string) error {
	game, err := s.gameRepo.FindByID(id)
	if err != nil {
		logger.Error("Game not found: %s", id)
		return errors.New("game not found")
	}

	if game == nil {
		logger.Error("Game not found: %s", id)
		return errors.New("game not found")
	}

	if game.UserID != userID {
		logger.Error("User %s does not own game %s", userID, id)
		return errors.New("access denied")
	}

	gameName := game.Name

	if err := s.gameRepo.Delete(id); err != nil {
		logger.Error("Failed to delete game %s: %v", id, err)
		return errors.New("failed to delete game")
	}

	// Notificación automática al eliminar juego
	go s.notificationService.SendNotificationToUser(
		userID,
		"Juego eliminado",
		"\""+gameName+"\" fue eliminado de tu biblioteca",
		"game_updates",
	)

	return nil
}
