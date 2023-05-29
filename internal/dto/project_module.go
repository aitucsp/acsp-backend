package dto

// CreateProjectModule DTO for Creating Project Module
type CreateProjectModule struct {
	Title        string `json:"title" form:"title" binding:"required" validate:"required"`
	Description  string `json:"description" form:"description" binding:"required" validate:"required"`
	ReferenceURL string `json:"reference_url" form:"reference_url"`
}

// UpdateProjectModule DTO for Updating Project Module
type UpdateProjectModule struct {
	Title        string `json:"title" form:"title" binding:"required" validate:"required"`
	Description  string `json:"description" form:"description" binding:"required" validate:"required"`
	ReferenceURL string `json:"reference_url" form:"reference_url"`
}
