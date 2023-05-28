package model

import (
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type User struct {
	ID        string         `json:"id" db:"id"`
	Email     string         `json:"email" db:"email"`
	Name      string         `json:"name" db:"name"`
	Password  string         `json:"-" db:"password"`
	CreatedAt string         `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt string         `json:"updated_at,omitempty" db:"updated_at"`
	IsAdmin   bool           `json:"-" db:"is_admin"`
	Roles     pq.StringArray `json:"-" db:"roles"`
	ImageURL  string         `json:"image_url" db:"image_url"`
	UserInfo  *UserDetails   `json:"user_details,omitempty" db:"user_details,omitempty"`
}

type UserDetails struct {
	ID             string `json:"id,omitempty" db:"id,omitempty"`
	UserID         string `json:"user_id,omitempty" db:"user_id,omitempty"`
	FirstName      string `json:"first_name" db:"first_name,omitempty"`
	LastName       string `json:"last_name" db:"last_name,omitempty"`
	PhoneNumber    string `json:"phone_number" db:"phone_number,omitempty"`
	Specialization string `json:"specialization" db:"specialization,omitempty"`
	UpdatedAt      string `json:"updated_at" db:"updated_at,omitempty"`
}
