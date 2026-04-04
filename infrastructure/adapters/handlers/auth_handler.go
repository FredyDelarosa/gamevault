package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gv/application/dto"
	"gv/application/mapper"
	"gv/core/logger"
	"gv/domain/ports/services"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register registra un nuevo usuario.
// @Summary Registrar usuario
// @Description Crea un nuevo usuario y retorna el token de autenticacion.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Datos de registro"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid register request: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, token, err := h.authService.Register(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.AuthResponse{
		Token: token,
		User:  mapper.ToUserResponse(user),
	})
}

// Login autentica un usuario.
// @Summary Iniciar sesion
// @Description Inicia sesion y retorna token JWT.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Credenciales"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid login request: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Token: token,
		User:  mapper.ToUserResponse(user),
	})
}
