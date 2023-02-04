package dto

// UpdateCard DTO for Updating Card
type UpdateCard struct {
	Position    string   `json:"position" binding:"required" validate:"required"`
	Skills      []string `json:"skills" binding:"required" validate:"required"`
	Description string   `json:"description" binding:"required" validate:"required"`
}

// CreateCard DTO for Creating Card
type CreateCard struct {
	Position    string   `json:"position" binding:"required" validate:"required"`
	Skills      []string `json:"skills" binding:"required" validate:"required"`
	Description string   `json:"description" binding:"required" validate:"required"`
}
