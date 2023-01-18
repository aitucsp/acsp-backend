package model

import (
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type User struct {
	ID        string         `json:"id" db:"id"`
	Email     string         `json:"email" db:"email"`
	Name      string         `json:"name" db:"name"`
	Password  string         `json:"password,omitempty" db:"password"`
	CreatedAt string         `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt string         `json:"updated_at,omitempty" db:"updated_at"`
	IsAdmin   bool           `json:"-" db:"is_admin"`
	Roles     pq.StringArray `json:"roles" db:"roles"`
}
