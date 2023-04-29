package model

type Contest struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"contest_name" binding:"required"`
	Description string `json:"description" db:"description" binding:"required"`
	Link        string `json:"link" db:"link" binding:"required"`
	StartDate   string `json:"start_date" db:"start_date" binding:"required"`
	EndDate     string `json:"end_date" db:"end_date" binding:"required"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}
