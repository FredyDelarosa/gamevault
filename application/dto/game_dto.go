package dto

type CreateGameRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	CoverImageURL string `json:"cover_image_url"`
	Status        string `json:"status" binding:"omitempty,oneof=NOW_PLAYING BACKLOG WISHLIST"`
}

type UpdateGameRequest struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	CoverImageURL *string `json:"cover_image_url"`
	Status        *string `json:"status" binding:"omitempty,oneof=NOW_PLAYING BACKLOG WISHLIST"`
	Completed     *bool   `json:"completed"`
}

type GameResponse struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	CoverImageURL string `json:"cover_image_url"`
	Status        string `json:"status"`
	Completed     bool   `json:"completed"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
