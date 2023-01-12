package model

type Comment struct {
	ID        int    `json:"id" db:"id"`
	UserID    int    `json:"user_id" db:"user_id"`
	ArticleID int    `json:"article_id" db:"article_id"`
	ParentID  int    `json:"parent_id" db:"parent_id,omitempty"`
	Text      string `json:"text" db:"text" binding:"required"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	Author    User   `json:"user" db:"user"`
}
