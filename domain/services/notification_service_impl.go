package services

import (
	"context"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
	"google.golang.org/api/option"

	"gv/core/config"
	"gv/core/logger"
	"gv/domain/models"
	"gv/domain/ports/repositories"
)

type NotificationServiceImpl struct {
	tokenRepo repositories.DeviceTokenRepository
	fcmClient *messaging.Client
}

func NewNotificationService(
	tokenRepo repositories.DeviceTokenRepository,
	cfg *config.Config,
) *NotificationServiceImpl {
	service := &NotificationServiceImpl{
		tokenRepo: tokenRepo,
	}

	// Inicializar Firebase Admin SDK
	if cfg.FirebaseCredentialsPath != "" {
		ctx := context.Background()
		app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(cfg.FirebaseCredentialsPath))
		if err != nil {
			logger.Error("Error inicializando Firebase: %v", err)
			return service
		}

		client, err := app.Messaging(ctx)
		if err != nil {
			logger.Error("Error obteniendo cliente FCM: %v", err)
			return service
		}

		service.fcmClient = client
		logger.Info("Firebase Cloud Messaging inicializado exitosamente")
	} else {
		logger.Info("Firebase credentials no configuradas, notificaciones push deshabilitadas")
	}

	return service
}

func (s *NotificationServiceImpl) SaveDeviceToken(userID, fcmToken string) error {
	token := &models.DeviceToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		FcmToken:  fcmToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.tokenRepo.SaveToken(token); err != nil {
		logger.Error("Error guardando token FCM: %v", err)
		return err
	}

	logger.Info("Token FCM guardado para usuario %s", userID)
	return nil
}

func (s *NotificationServiceImpl) SendNotificationToUser(userID, title, body, channel string) error {
	if s.fcmClient == nil {
		logger.Info("FCM client no disponible, omitiendo notificación")
		return nil
	}

	// Obtener todos los tokens del usuario
	tokens, err := s.tokenRepo.FindByUserID(userID)
	if err != nil {
		logger.Error("Error obteniendo tokens para usuario %s: %v", userID, err)
		return err
	}

	if len(tokens) == 0 {
		logger.Info("No hay dispositivos registrados para usuario %s", userID)
		return nil
	}

	ctx := context.Background()

	// Enviar notificación a cada dispositivo del usuario
	for _, token := range tokens {
		message := &messaging.Message{
			Token: token.FcmToken,
			Notification: &messaging.Notification{
				Title: title,
				Body:  body,
			},
			Data: map[string]string{
				"channel": channel,
				"title":   title,
				"body":    body,
			},
			Android: &messaging.AndroidConfig{
				Priority: "high",
				Notification: &messaging.AndroidNotification{
					ChannelID: channel,
				},
			},
		}

		_, err := s.fcmClient.Send(ctx, message)
		if err != nil {
			logger.Error("Error enviando notificación a token %s: %v", token.FcmToken[:20], err)
			// Si el token es inválido, eliminarlo
			if messaging.IsUnregistered(err) {
				logger.Info("Token inválido, eliminando: %s", token.FcmToken[:20])
				s.tokenRepo.DeleteByToken(token.FcmToken)
			}
		} else {
			logger.Info("Notificación enviada exitosamente a usuario %s", userID)
		}
	}

	return nil
}