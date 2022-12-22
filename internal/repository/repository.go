package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"acsp/internal/model"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetUser(ctx context.Context, username, password string) (*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(ctx context.Context) (*[]model.User, error)
}

type Articles interface {
	Create(ctx context.Context, article model.Article) (int64, error)
	Update(ctx context.Context, article model.Article) (int64, error)
	Delete(ctx context.Context, userID int, articleID int) (int64, error)
	GetAllByUserId(ctx context.Context, userID int) ([]model.Article, error)
}

type Repository struct {
	Authorization
	Articles
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Articles:      NewArticlesPostgres(db),
	}
}
