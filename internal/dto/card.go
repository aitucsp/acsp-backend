package dto

// UpdateCard DTO for Updating Card
type UpdateCard struct {
	Position    string   `json:"position" form:"position" binding:"required" validate:"required"`
	Skills      []string `json:"skills" form:"skills" binding:"required" validate:"required"`
	Description string   `json:"description" form:"description" binding:"required" validate:"required"`
}

// CreateCard DTO for Creating Card
type CreateCard struct {
	Position    string   `json:"position" form:"position" binding:"required" validate:"required"`
	Skills      []string `json:"skills" form:"skills" binding:"required" validate:"required"`
	Description string   `json:"description" form:"description" binding:"required" validate:"required"`
}

// AnswerInvitation DTO for Creating Invitation Card
type AnswerInvitation struct {
	Feedback string `json:"feedback" form:"feedback" binding:"required" validate:"required"`
}
