package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"gv/application/dto"
	"gv/core/logger"
	"gv/domain/ports/services"
)

type AuthMiddleware struct {
	authService services.AuthService
}

func NewAuthMiddleware(authService services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Error("Missing authorization header")
			c.JSON(401, dto.ErrorResponse{Error: "missing authorization header"})
			c.Abort()
			return
		}

		// Formato esperado: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Error("Invalid authorization header format")
			c.JSON(401, dto.ErrorResponse{Error: "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		userID, err := m.authService.ValidateToken(token)
		if err != nil {
			logger.Error("Invalid token: %v", err)
			c.JSON(401, dto.ErrorResponse{Error: "invalid or expired token"})
			c.Abort()
			return
		}

		// Guardar userID en el contexto para uso posterior
		c.Set("userID", userID)
		c.Next()
	}
}

// GetUserID obtiene el userID del contexto
func GetUserID(c *gin.Context) string {
	userID, exists := c.Get("userID")
	if !exists {
		return ""
	}
	return userID.(string)
}
