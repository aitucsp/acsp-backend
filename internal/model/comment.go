package model

import "database/sql"

type Comment struct {
	ID        int           `json:"id" db:"id"`
	UserID    int           `json:"user_id" db:"user_id"`
	ArticleID int           `json:"article_id" db:"article_id"`
	ParentID  sql.NullInt64 `json:"parent_id,omitempty" db:"parent_id,omitempty"`
	Text      string        `json:"text" db:"text" binding:"required"`
	CreatedAt string        `json:"created_at" db:"created_at"`
	UpdatedAt string        `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	Author    *User         `json:"user" db:"-"`
}
