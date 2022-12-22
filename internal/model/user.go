package model

import (
	_ "github.com/lib/pq"
)

type User struct {
	ID        string `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Name      string `json:"name" db:"name"`
	Password  string `json:"password" db:"password"`
	CreatedAt string `json:"-" db:"created_at"`
	UpdatedAt string `json:"-" db:"updated_at"`
}
