package service

import (
	"context"

	_ "github.com/golang/mock/gomock"

	"gorestapi/internal/dto"
	"gorestapi/internal/model"
	"gorestapi/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(ctx context.Context, dto dto.CreateUser) (string, error)
	GenerateToken(ctx context.Context, email, password string) (string, error)
	ParseToken(token string) (string, error)
}

type Projects interface {
	Create(ctx context.Context, userID string, dto dto.CreateProject) (string, error)
	Update(ctx context.Context, userID string, projectId string, newTitle string) error
	Delete(ctx context.Context, userID string, projectId string) error
	GetByTitle(ctx context.Context, userID string, title string) (model.Project, error)
	GetAll(ctx context.Context) ([]model.Project, error)
}

type Service struct {
	Authorization
	Projects
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Projects:      NewProjectsService(repo.Projects),
	}
}
