package dto

type CreateUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required" validate:"required"`
	Email    string `json:"email" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
}
