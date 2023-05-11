package dto

// UpdateComment DTO for Updating Comment
type UpdateComment struct {
	Topic       string `json:"topic" binding:"required" validate:"required"`
	Description string `json:"description" binding:"required" validate:"required"`
}

// CreateComment DTO for Commenting Article
type CreateComment struct {
	Text string `json:"text" binding:"required" validate:"required"`
}

// ReplyToComment DTO for Replying to a Comment of an Article
type ReplyToComment struct {
	Text string `json:"text" binding:"required" validate:"required"`
}
