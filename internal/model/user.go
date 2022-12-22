package model

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name" binding:"required"`
	Email    string `json:"email" bson:"email" binding:"required" validate:"email"`
	Password string `json:"password" bson:"password" binding:"required" validate:"min=8"`
}
