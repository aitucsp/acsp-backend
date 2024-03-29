package model

type Material struct {
	ID          int    `json:"id" db:"id"`
	UserID      int    `json:"user_id" db:"user_id"`
	Topic       string `json:"topic" db:"topic" binding:"required"`
	Description string `json:"description" db:"description" binding:"required"`
	Upvote      int    `json:"upvote" db:"upvote"`
	Downvote    int    `json:"downvote" db:"downvote"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
	Author      *User  `json:"-" db:"-"`
}
