package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"gv/core/config"
	"gv/core/logger"
	"gv/domain/models"
	"gv/domain/ports/repositories"
	"gv/domain/ports/services"
)

type AuthServiceImpl struct {
	userRepo            repositories.UserRepository
	cfg                 *config.Config
	notificationService services.NotificationService
}

func NewAuthService(userRepo repositories.UserRepository, cfg *config.Config, notificationService services.NotificationService) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo:            userRepo,
		cfg:                 cfg,
		notificationService: notificationService,
	}
}

func (s *AuthServiceImpl) Register(email, password, firstName, lastName string) (*models.User, string, error) {
	// Verificar si el usuario ya existe
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		logger.Error("User already exists: %s", email)
		return nil, "", errors.New("user already exists")
	}

	// Hashear contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password: %v", err)
		return nil, "", errors.New("failed to process password")
	}

	// Crear usuario
	user := &models.User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("Failed to create user: %v", err)
		return nil, "", errors.New("failed to create user")
	}

	// Generar token JWT
	token, err := s.generateToken(user.ID)
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		return nil, "", errors.New("failed to generate authentication token")
	}

	// Notificación de bienvenida
	go s.notificationService.SendNotificationToUser(
		user.ID,
		"¡Bienvenido a GameVault!",
		"Hola "+firstName+", tu cuenta fue creada exitosamente",
		"game_updates",
	)

	return user, token, nil
}

func (s *AuthServiceImpl) Login(email, password string) (*models.User, string, error) {
	// Buscar usuario por email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		logger.Error("User not found: %s", email)
		return nil, "", errors.New("invalid credentials")
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logger.Error("Invalid password for user: %s", email)
		return nil, "", errors.New("invalid credentials")
	}

	// Generar token JWT
	token, err := s.generateToken(user.ID)
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		return nil, "", errors.New("failed to generate authentication token")
	}

	return user, token, nil
}

func (s *AuthServiceImpl) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil {
		logger.Error("Failed to parse token: %v", err)
		return "", errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid token claims")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}

func (s *AuthServiceImpl) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expira en 24 horas
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
