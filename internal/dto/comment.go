package dto

// UpdateComment DTO for Updating Comment
type UpdateComment struct {
	Topic       string `json:"topic" form:"topic" binding:"required" validate:"required"`
	Description string `json:"description" form:"description" binding:"required" validate:"required"`
}

// CreateComment DTO for Commenting Article
type CreateComment struct {
	Text string `json:"text" form:"text" binding:"required" validate:"required"`
}

// ReplyToComment DTO for Replying to a Comment of an Article
type ReplyToComment struct {
	Text string `json:"text" form:"text" binding:"required" validate:"required"`
}
