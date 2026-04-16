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

type PostServiceImpl struct {
	postRepo            repositories.PostRepository
	reactionRepo        repositories.PostReactionRepository
	gameRepo            repositories.GameRepository
	userRepo            repositories.UserRepository
	notificationService services.NotificationService
}

func NewPostService(
	postRepo repositories.PostRepository,
	reactionRepo repositories.PostReactionRepository,
	gameRepo repositories.GameRepository,
	userRepo repositories.UserRepository,
	notificationService services.NotificationService,
) *PostServiceImpl {
	return &PostServiceImpl{
		postRepo:            postRepo,
		reactionRepo:        reactionRepo,
		gameRepo:            gameRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
	}
}

func (s *PostServiceImpl) CreatePost(userID, gameName, title, content string, postType models.PostType) (*models.Post, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if gameName == "" {
		return nil, errors.New("game name is required")
	}

	if postType == "" {
		postType = models.PostTypeDiscussion
	}

	post := &models.Post{
		ID:             uuid.New().String(),
		UserID:         userID,
		GameName:       gameName,
		Title:          title,
		Content:        content,
		PostType:       postType,
		ReactionsCount: 0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.postRepo.Create(post); err != nil {
		logger.Error("Failed to create post: %v", err)
		return nil, errors.New("failed to create post")
	}

	// Notificar a usuarios que tienen ese juego en su biblioteca
	go s.notifyUsersWithGame(post, userID)

	return post, nil
}

func (s *PostServiceImpl) notifyUsersWithGame(post *models.Post, authorID string) {
	// Buscar todos los juegos con ese nombre (case-insensitive)
	users, err := s.findUsersWithGame(post.GameName)
	if err != nil {
		logger.Error("Error buscando usuarios con juego %s: %v", post.GameName, err)
		return
	}

	// Obtener nombre del autor
	author, err := s.userRepo.FindByID(authorID)
	authorName := "Alguien"
	if err == nil && author != nil {
		authorName = author.FirstName
	}

	// Construir título según tipo de post
	typeLabel := ""
	switch post.PostType {
	case models.PostTypeTip:
		typeLabel = "dio un tip"
	case models.PostTypeReview:
		typeLabel = "publicó una reseña"
	case models.PostTypeQuestion:
		typeLabel = "hizo una pregunta"
	case models.PostTypeNews:
		typeLabel = "compartió una noticia"
	default:
		typeLabel = "publicó"
	}

	title := authorName + " " + typeLabel + " sobre " + post.GameName
	body := post.Title

	// Notificar a cada usuario (excepto al autor)
	notified := make(map[string]bool)
	for _, userID := range users {
		if userID == authorID || notified[userID] {
			continue
		}
		notified[userID] = true
		s.notificationService.SendNotificationToUser(userID, title, body, "game_updates")
	}
}

func (s *PostServiceImpl) findUsersWithGame(gameName string) ([]string, error) {
	// Buscar todos los juegos cuyo nombre haga match (case-insensitive)
	// Usamos GORM directamente para eficiencia
	// Estrategia: buscar todos los juegos y filtrar
	// En producción convendría un query más eficiente
	allGames, err := s.findAllGamesMatchingName(gameName)
	if err != nil {
		return nil, err
	}

	userIDs := []string{}
	seen := make(map[string]bool)
	for _, game := range allGames {
		if !seen[game.UserID] {
			seen[game.UserID] = true
			userIDs = append(userIDs, game.UserID)
		}
	}
	return userIDs, nil
}

func (s *PostServiceImpl) findAllGamesMatchingName(gameName string) ([]models.Game, error) {
	// Hack: necesitamos acceso directo al DB para esta query
	// Agregaremos un método al GameRepository
	return s.gameRepo.FindByNameMatch(gameName)
}

func (s *PostServiceImpl) GetAllPosts(limit, offset int) ([]models.Post, error) {
	if limit <= 0 {
		limit = 20
	}
	posts, err := s.postRepo.FindAll(limit, offset)
	if err != nil {
		logger.Error("Failed to get posts: %v", err)
		return nil, errors.New("failed to retrieve posts")
	}
	return posts, nil
}

func (s *PostServiceImpl) GetPostsForUserGames(userID string, limit, offset int) ([]models.Post, error) {
	if limit <= 0 {
		limit = 20
	}

	// Obtener juegos del usuario
	games, err := s.gameRepo.FindByUserID(userID, nil)
	if err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return []models.Post{}, nil
	}

	gameNames := []string{}
	for _, g := range games {
		gameNames = append(gameNames, g.Name)
	}

	posts, err := s.postRepo.FindByGameNames(gameNames, limit, offset)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostServiceImpl) GetPostByID(id string) (*models.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("post not found")
	}
	if post == nil {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (s *PostServiceImpl) DeletePost(id, userID string) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil || post == nil {
		return errors.New("post not found")
	}

	if post.UserID != userID {
		return errors.New("access denied")
	}

	return s.postRepo.Delete(id)
}

func (s *PostServiceImpl) ToggleReaction(postID, userID string) (bool, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil || post == nil {
		return false, errors.New("post not found")
	}

	exists, err := s.reactionRepo.Exists(postID, userID)
	if err != nil {
		return false, err
	}

	if exists {
		// Ya reaccionó → quitar reacción
		if err := s.reactionRepo.Delete(postID, userID); err != nil {
			return false, err
		}
		s.postRepo.DecrementReactions(postID)
		return false, nil
	}

	// No había reaccionado → agregar reacción
	reaction := &models.PostReaction{
		ID:        uuid.New().String(),
		PostID:    postID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := s.reactionRepo.Create(reaction); err != nil {
		return false, err
	}
	s.postRepo.IncrementReactions(postID)

	// Notificar al autor del post (si no es el mismo que reacciona)
	if post.UserID != userID {
		go s.notifyPostAuthor(post, userID)
	}

	return true, nil
}

func (s *PostServiceImpl) notifyPostAuthor(post *models.Post, reactorID string) {
	reactor, err := s.userRepo.FindByID(reactorID)
	reactorName := "Alguien"
	if err == nil && reactor != nil {
		reactorName = reactor.FirstName
	}

	title := "¡Tu publicación recibió una reacción!"
	body := reactorName + " recomendó tu publicación sobre " + post.GameName

	s.notificationService.SendNotificationToUser(post.UserID, title, body, "game_updates")
}

func (s *PostServiceImpl) HasUserReacted(postID, userID string) (bool, error) {
	return s.reactionRepo.Exists(postID, userID)
}
