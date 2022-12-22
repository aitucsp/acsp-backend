package model

type Article struct {
	ID          int    `json:"id" db:"id"`
	UserID      int    `json:"-" db:"user_id"`
	Topic       string `json:"topic" db:"topic" binding:"required"`
	Description string `json:"description" db:"description" binding:"required"`
	CreatedAt   string `json:"-" db:"created_at"`
	UpdatedAt   string `json:"-" db:"updated_at"`
	Author      *User  `json:"-" db:"-"`
}
