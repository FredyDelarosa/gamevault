package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gv/application/dto"
	"gv/application/mapper"
	"gv/core/logger"
	"gv/domain/models"
	"gv/domain/ports/services"
	"gv/infrastructure/middleware"
)

type GameHandler struct {
	gameService services.GameService
}

func NewGameHandler(gameService services.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// CreateGame crea un nuevo juego para el usuario autenticado.
// @Summary Crear juego
// @Description Crea un juego asociado al usuario autenticado.
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param request body dto.CreateGameRequest true "Datos del juego"
// @Success 201 {object} dto.GameResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/games [post]
func (h *GameHandler) CreateGame(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid create game request: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Convertir status string a enum
	status := models.GameStatus(req.Status)
	if req.Status == "" {
		status = models.StatusWishlist
	}

	game, err := h.gameService.CreateGame(
		userID,
		req.Name,
		req.Description,
		req.CoverImageURL,
		status,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapper.ToGameResponse(game))
}

// GetGames lista juegos del usuario autenticado.
// @Summary Listar juegos
// @Description Retorna el listado de juegos del usuario autenticado.
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param status query string false "Filtrar por estado" Enums(NOW_PLAYING, BACKLOG, WISHLIST)
// @Success 200 {array} dto.GameResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/games [get]
func (h *GameHandler) GetGames(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	// Filtro opcional por status
	statusParam := c.Query("status")
	var status *models.GameStatus
	if statusParam != "" {
		s := models.GameStatus(statusParam)
		status = &s
	}

	games, err := h.gameService.GetGamesByUser(userID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapper.ToGameResponseList(games))
}

// GetGameByID obtiene un juego por ID.
// @Summary Obtener juego por ID
// @Description Retorna un juego por su ID si pertenece al usuario autenticado.
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del juego"
// @Success 200 {object} dto.GameResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/games/{id} [get]
func (h *GameHandler) GetGameByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "game id is required"})
		return
	}

	game, err := h.gameService.GetGameByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapper.ToGameResponse(game))
}

// UpdateGame actualiza un juego existente.
// @Summary Actualizar juego
// @Description Actualiza campos de un juego del usuario autenticado.
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del juego"
// @Param request body dto.UpdateGameRequest true "Campos a actualizar"
// @Success 200 {object} dto.GameResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/games/{id} [put]
func (h *GameHandler) UpdateGame(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "game id is required"})
		return
	}

	var req dto.UpdateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid update game request: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Convertir status si viene
	var status *models.GameStatus
	if req.Status != nil {
		s := models.GameStatus(*req.Status)
		status = &s
	}

	game, err := h.gameService.UpdateGame(
		id,
		userID,
		req.Name,
		req.Description,
		req.CoverImageURL,
		status,
		req.Completed,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapper.ToGameResponse(game))
}

// DeleteGame elimina un juego existente.
// @Summary Eliminar juego
// @Description Elimina un juego del usuario autenticado.
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del juego"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/games/{id} [delete]
func (h *GameHandler) DeleteGame(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "game id is required"})
		return
	}

	err := h.gameService.DeleteGame(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "game deleted successfully"})
}
