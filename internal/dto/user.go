package dto

// CreateUser DTO for Creating a User
type CreateUser struct {
	ID       int    `json:"-"`
	Name     string `json:"name" form:"name" binding:"required" validate:"required"`
	Email    string `json:"email" form:"email" binding:"required" validate:"required,email,endswith=@astanait.edu.kz"`
	Password string `json:"password" form:"password" binding:"required" validate:"required"`
}

// UpdateUser DTO for Updating a User
type UpdateUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name" form:"name" binding:"required" validate:"required"`
	Email    string `json:"email" form:"email" binding:"required" validate:"required,email,endswith=@astanait.edu.kz"`
	Password string `json:"password" form:"password" binding:"required" validate:"required"`
}
