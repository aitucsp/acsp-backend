package model

import (
	_ "github.com/lib/pq"
)

type Project struct {
	Id       *string `json:"id" bson:"_id,omitempty"`
	Title    string  `json:"title"  bson:"title" validate:"min=2"`
	Archived bool    `json:"archived" bson:"archived"`
	UserID   string  `json:"-" bson:"user_id"`
}

func (p *Project) Archive() {
	p.Archived = true
}

func (p *Project) Restore() {
	p.Archived = false
}
