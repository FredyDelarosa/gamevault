package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gv/application/dto"
	"gv/core/logger"
	"gv/domain/models"
	"gv/domain/ports/repositories"
	"gv/domain/ports/services"
	"gv/infrastructure/middleware"
)

type PostHandler struct {
	postService services.PostService
	userRepo    repositories.UserRepository
}

func NewPostHandler(postService services.PostService, userRepo repositories.UserRepository) *PostHandler {
	return &PostHandler{
		postService: postService,
		userRepo:    userRepo,
	}
}

func (h *PostHandler) toResponse(post *models.Post, currentUserID string) dto.PostResponse {
	authorName := "Usuario"
	if user, err := h.userRepo.FindByID(post.UserID); err == nil && user != nil {
		authorName = user.FirstName + " " + user.LastName
	}

	hasReacted, _ := h.postService.HasUserReacted(post.ID, currentUserID)

	return dto.PostResponse{
		ID:             post.ID,
		UserID:         post.UserID,
		AuthorName:     authorName,
		GameName:       post.GameName,
		Title:          post.Title,
		Content:        post.Content,
		PostType:       string(post.PostType),
		ReactionsCount: post.ReactionsCount,
		HasReacted:     hasReacted,
		CreatedAt:      post.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      post.UpdatedAt.Format(time.RFC3339),
	}
}

// CreatePost crea una nueva publicación.
// @Summary Crear publicación
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePostRequest true "Datos del post"
// @Success 201 {object} dto.PostResponse
// @Router /api/posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Invalid post request: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	postType := models.PostType(req.PostType)
	if req.PostType == "" {
		postType = models.PostTypeDiscussion
	}

	post, err := h.postService.CreatePost(userID, req.GameName, req.Title, req.Content, postType)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.toResponse(post, userID))
}

// GetAllPosts obtiene todos los posts (feed general).
// @Summary Listar todos los posts
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Límite de resultados"
// @Param offset query int false "Offset para paginación"
// @Success 200 {array} dto.PostResponse
// @Router /api/posts [get]
func (h *PostHandler) GetAllPosts(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	posts, err := h.postService.GetAllPosts(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	responses := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = h.toResponse(&post, userID)
	}

	c.JSON(http.StatusOK, responses)
}

// GetPostsForMyGames obtiene posts de juegos en la biblioteca del usuario.
// @Summary Feed filtrado por mis juegos
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.PostResponse
// @Router /api/posts/my-games [get]
func (h *PostHandler) GetPostsForMyGames(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	posts, err := h.postService.GetPostsForUserGames(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	responses := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = h.toResponse(&post, userID)
	}

	c.JSON(http.StatusOK, responses)
}

// GetPostByID obtiene un post específico.
// @Summary Obtener post por ID
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID del post"
// @Success 200 {object} dto.PostResponse
// @Router /api/posts/{id} [get]
func (h *PostHandler) GetPostByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id := c.Param("id")
	post, err := h.postService.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.toResponse(post, userID))
}

// DeletePost elimina un post propio.
// @Summary Eliminar post
// @Tags posts
// @Security BearerAuth
// @Param id path string true "ID del post"
// @Success 200 {object} dto.SuccessResponse
// @Router /api/posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id := c.Param("id")
	if err := h.postService.DeletePost(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "post deleted successfully"})
}

// ToggleReaction alterna la reacción (like) a un post.
// @Summary Reaccionar / quitar reacción
// @Tags posts
// @Security BearerAuth
// @Param id path string true "ID del post"
// @Success 200 {object} dto.SuccessResponse
// @Router /api/posts/{id}/react [post]
func (h *PostHandler) ToggleReaction(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id := c.Param("id")
	reacted, err := h.postService.ToggleReaction(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	message := "reaction removed"
	if reacted {
		message = "reaction added"
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: message})
}
