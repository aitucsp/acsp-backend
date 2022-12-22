package service

import (
	"context"

	"gorestapi/internal/dto"
	"gorestapi/internal/model"
	"gorestapi/internal/repository"
)

type ProjectsService struct {
	repo repository.Projects
}

func NewProjectsService(repo repository.Projects) *ProjectsService {
	return &ProjectsService{repo: repo}
}

func (s *ProjectsService) Create(ctx context.Context, userID string, dto dto.CreateProject) (string, error) {
	project := model.Project{
		Title:    dto.Title,
		UserID:   userID,
		Archived: false,
	}

	return s.repo.Create(ctx, project)
}

func (s *ProjectsService) Update(ctx context.Context, userID string, projectId string, newTitle string) error {
	return s.repo.Update(ctx, userID, projectId, newTitle)
}

func (s *ProjectsService) GetAll(ctx context.Context) ([]model.Project, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProjectsService) GetByTitle(ctx context.Context, userID string, title string) (model.Project, error) {
	return s.repo.GetByTitle(ctx, userID, title)
}

func (s *ProjectsService) Delete(ctx context.Context, userID string, projectId string) error {
	return s.repo.Delete(ctx, userID, projectId)
}
