package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gv/application/dto"
	"gv/core/logger"
	"gv/domain/ports/services"
	"gv/infrastructure/middleware"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// SaveFcmToken guarda el token FCM del dispositivo.
// @Summary Guardar token FCM
// @Description Guarda el token FCM del dispositivo para recibir push notifications.
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param request body dto.SaveFcmTokenRequest true "Token FCM"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/notifications/token [post]
func (h *NotificationHandler) SaveFcmToken(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.SaveFcmTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid FCM token request: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.notificationService.SaveDeviceToken(userID, req.FcmToken); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to save token"})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "token saved successfully"})
}

// SendNotification envía una notificación push a un usuario.
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid notification request: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	channel := req.Channel
	if channel == "" {
		channel = "game_updates"
	}

	if err := h.notificationService.SendNotificationToUser(userID, req.Title, req.Body, channel); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to send notification"})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "notification sent successfully"})
}
