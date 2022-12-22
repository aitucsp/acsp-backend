package dto

// UpdateProject DTO for Updating Project
type UpdateProject struct {
	NewTitle *string `json:"new_title" validate:"required"`
}

// CreateProject DTO for Creating Project
type CreateProject struct {
	Title string `json:"title" binding:"required" validate:"required"`
}
