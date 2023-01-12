package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"acsp/internal/model"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, username, password string) (*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(ctx context.Context) (*[]model.User, error)
}

type Roles interface {
	CreateRole(ctx context.Context, name string) error
	UpdateRole(ctx context.Context, roleID int, newName string) error
	DeleteRole(ctx context.Context, roleID int) error
	SaveUserRole(ctx context.Context, userID, roleID int) error
	DeleteUserRole(ctx context.Context, userID, roleID int) error
	GetUserRoles(ctx context.Context, userID int) ([]model.Role, error)
}

type Articles interface {
	Create(ctx context.Context, article model.Article) error
	Update(ctx context.Context, article model.Article) error
	Delete(ctx context.Context, userID int, articleID int) error
	GetAllByUserId(ctx context.Context, userID int) ([]model.Article, error)
	GetArticleByIDAndUserID(ctx context.Context, articleID, userID int) (model.Article, error)
	CreateComment(ctx context.Context, articleID, userID int, comment model.Comment) error
	GetCommentsByArticleID(ctx context.Context, articleID int) ([]model.Comment, error)
	ReplyToComment(ctx context.Context, articleID, userID, parentCommentID int, comment model.Comment) error
	GetRepliesByArticleIDAndCommentID(ctx context.Context, articleID, userID, parentCommentID int) ([]model.Comment, error)
}

type Repository struct {
	Authorization
	Roles
	Articles
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Roles:         NewRolesPostgres(db),
		Articles:      NewArticlesPostgres(db),
	}
}
