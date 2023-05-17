package dto

// CreateDiscipline is a DTO for creating a discipline
type CreateDiscipline struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateDiscipline is a DTO for updating a discipline
type UpdateDiscipline struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}
