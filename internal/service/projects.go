package service

import (
	"context"

	"acsp/internal/dto"
	"acsp/internal/model"
	"acsp/internal/repository"
)

type ProjectsService struct {
	repo repository.Projects
}

func NewProjectsService(repo repository.Projects) *ProjectsService {
	return &ProjectsService{repo: repo}
}

func (p ProjectsService) Create(ctx context.Context, project dto.CreateProject) error {
	// TODO implement me
	panic("implement me")
}

func (p ProjectsService) Update(ctx context.Context, project dto.UpdateProject) error {
	// TODO implement me
	panic("implement me")
}

func (p ProjectsService) Delete(ctx context.Context, projectID int) error {
	// TODO implement me
	panic("implement me")
}

func (p ProjectsService) GetAll(ctx context.Context) ([]model.Project, error) {
	// TODO implement me
	panic("implement me")
}

func (p ProjectsService) GetByID(ctx context.Context, projectID int) (model.Project, error) {
	// TODO implement me
	panic("implement me")
}
