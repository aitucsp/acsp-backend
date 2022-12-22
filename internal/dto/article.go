package dto

import "acsp/internal/model"

// UpdateArticle DTO for Updating Project
type UpdateArticle struct {
	ID          int        `json:"id" binding:"required"`
	Topic       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"required"`
	Author      model.User `json:"-"`
}

// CreateArticle DTO for Creating Project
type CreateArticle struct {
	Topic       string `json:"topic" binding:"required"`
	Description string `json:"description" binding:"required"`
}
