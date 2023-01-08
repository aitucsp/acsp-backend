package dto

import "acsp/internal/model"

// UpdateArticle DTO for Updating Article
type UpdateArticle struct {
	Topic       string     `json:"topic" binding:"required" validate:"required"`
	Description string     `json:"description" binding:"required" validate:"required"`
	Author      model.User `json:"-"`
}

// CreateArticle DTO for Creating Article
type CreateArticle struct {
	Topic       string `json:"topic" binding:"required" validate:"required"`
	Description string `json:"description" binding:"required" validate:"required"`
}
