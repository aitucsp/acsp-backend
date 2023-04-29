package model

import "github.com/lib/pq"

type Card struct {
	ID          int            `json:"id" db:"id"`
	UserID      int            `json:"user_id" db:"user_id"`
	Position    string         `json:"position" db:"position"`
	Skills      pq.StringArray `json:"skills" db:"skills"`
	Description string         `json:"description" db:"description"`
	CreatedAt   string         `json:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt   string         `json:"updated_at,omitempty" db:"updated_at,omitempty"`
	Author      User           `json:"user" db:"user,omitempty"`
}

type InvitationCard struct {
	Card      *Card  `json:"card" db:"-"`
	InviterID int    `json:"inviter_id" db:"inviter_id"`
	Status    string `json:"status" db:"status"`
}
