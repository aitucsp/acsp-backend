package service

import (
	"context"

	_ "github.com/golang/mock/gomock"

	"acsp/internal/dto"
	"acsp/internal/model"
	"acsp/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(ctx context.Context, dto dto.CreateUser) error
	GenerateToken(ctx context.Context, email, password string) (string, error)
	ParseToken(token string) (string, error)
}

type Roles interface {
}

type Articles interface {
	Create(ctx context.Context, userID string, dto dto.CreateArticle) error
	GetAll(ctx context.Context, userID string) ([]model.Article, error)
	GetByID(ctx context.Context, articleID, userID string) (model.Article, error)
	Update(ctx context.Context, articleID string, userID string, article dto.UpdateArticle) error
	Delete(ctx context.Context, userID string, projectId string) error
	CommentByID(ctx context.Context, articleID, userID string, comment dto.CreateComment) error
	GetCommentsByArticleID(ctx context.Context, articleID string) ([]model.Comment, error)
	ReplyToCommentByArticleIDAndCommentID(
		ctx context.Context, articleID string, userID string, parentCommentID string, comment dto.ReplyToComment) error
	GetRepliesByArticleIDAndCommentID(
		ctx context.Context, articleID, userID, commentID string) ([]model.Comment, error)
}

type Service struct {
	Authorization
	Articles
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Articles:      NewArticlesService(repo.Articles, repo.Authorization),
	}
}
