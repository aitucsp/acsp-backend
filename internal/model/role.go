package model

import (
	_ "github.com/lib/pq"
)

type Role struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
