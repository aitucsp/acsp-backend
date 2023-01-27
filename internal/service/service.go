package service

import (
	"context"

	"github.com/go-redis/redis/v9"
	_ "github.com/golang/mock/gomock"

	"acsp/internal/config"
	"acsp/internal/dto"
	"acsp/internal/model"
	"acsp/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Service struct {
	Authorization
	Articles
	Roles
}

type Authorization interface {
	CreateUser(ctx context.Context, dto dto.CreateUser) error
	GenerateTokenPair(ctx context.Context, email, password string) (*TokenDetails, error)
	ParseToken(token string) (string, error)
	SaveTokenPair(ctx context.Context, userID string, details *TokenDetails) error
}

type Roles interface {
	CreateRole(ctx context.Context, userID, name string) error
	UpdateRole(ctx context.Context, userID, roleID, newName string) error
	DeleteRole(ctx context.Context, userID, roleID string) error
	SaveUserRole(ctx context.Context, userID, roleID string) error
	DeleteUserRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]model.Role, error)
}

type Articles interface {
	Create(ctx context.Context, userID string, dto dto.CreateArticle) error
	GetAll(ctx context.Context, userID string) (*[]model.Article, error)
	GetByID(ctx context.Context, articleID, userID string) (*model.Article, error)
	Update(ctx context.Context, articleID string, userID string, article dto.UpdateArticle) error
	Delete(ctx context.Context, userID string, projectId string) error
	CommentByID(ctx context.Context, articleID, userID string, comment dto.CreateComment) error
	GetCommentsByArticleID(ctx context.Context, articleID string) ([]model.Comment, error)
	ReplyToCommentByArticleIDAndCommentID(
		ctx context.Context, articleID string, userID string, parentCommentID string, comment dto.ReplyToComment) error
	GetRepliesByArticleIDAndCommentID(
		ctx context.Context, articleID, userID, commentID string) (*[]model.Comment, error)
}

func NewService(repo *repository.Repository, r *redis.Client, c config.AuthConfig) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization, repo.Roles, r, c),
		Articles:      NewArticlesService(repo.Articles, repo.Authorization),
		Roles:         NewRolesService(repo.Roles, repo.Authorization),
	}
}
