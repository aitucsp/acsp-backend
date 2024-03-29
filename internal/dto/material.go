package dto

import "acsp/internal/model"

// UpdateMaterial DTO for Updating Article
type UpdateMaterial struct {
	Topic       string     `json:"topic" form:"topic" binding:"required" validate:"required"`
	Description string     `json:"description" form:"description" binding:"required" validate:"required"`
	Author      model.User `json:"-"`
}

// CreateMaterial DTO for Creating Article
type CreateMaterial struct {
	Topic       string `json:"topic" form:"topic" binding:"required" validate:"required"`
	Description string `json:"description" form:"description" binding:"required" validate:"required"`
}
