package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"gorestapi/internal/model"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (string, error)
	GetUser(ctx context.Context, username, password string) (*model.User, error)
	GetById(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(ctx context.Context) (*[]model.User, error)
}

type Projects interface {
	Create(ctx context.Context, project model.Project) (string, error)
	Update(ctx context.Context, userId string, projectId string, newTitle string) error
	Delete(ctx context.Context, userId string, projectId string) error
	GetAll(ctx context.Context) ([]model.Project, error)
	GetByTitle(ctx context.Context, userId string, title string) (model.Project, error)
}

type Repository struct {
	Authorization
	Projects
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Authorization: NewAuthMongo(db),
		Projects:      NewProjectsDB(db),
	}
}
