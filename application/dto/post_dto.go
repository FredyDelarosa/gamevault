package dto

type CreatePostRequest struct {
	GameName string `json:"game_name" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
	PostType string `json:"post_type" binding:"omitempty,oneof=TIP DISCUSSION REVIEW QUESTION NEWS"`
}

type PostResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	AuthorName     string `json:"author_name"`
	GameName       string `json:"game_name"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	PostType       string `json:"post_type"`
	ReactionsCount int    `json:"reactions_count"`
	HasReacted     bool   `json:"has_reacted"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
