package dto

// CreateProjectModule DTO for Creating Project Module
type CreateProjectModule struct {
	Title string `json:"title" binding:"required" validate:"required"`
}

// UpdateProjectModule DTO for Updating Project Module
type UpdateProjectModule struct {
	Title string `json:"title" binding:"required" validate:"required"`
}
