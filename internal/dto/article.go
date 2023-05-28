package dto

import (
	"mime/multipart"

	"acsp/internal/model"
)

// CreateArticle DTO for Creating Article
type CreateArticle struct {
	Topic       string                `json:"topic" form:"topic" binding:"required" validate:"required"`
	Description string                `json:"description" form:"description" binding:"required" validate:"required"`
	Image       *multipart.FileHeader `form:"-" json:"-"`
}

// UpdateArticle DTO for Updating Article
type UpdateArticle struct {
	Topic       string     `json:"topic" form:"topic" binding:"required" validate:"required"`
	Description string     `json:"description" form:"description" binding:"required" validate:"required"`
	Author      model.User `json:"-"`
}
