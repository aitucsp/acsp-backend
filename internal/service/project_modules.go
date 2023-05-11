package service

import (
	"context"

	"github.com/pkg/errors"

	"acsp/internal/dto"
	"acsp/internal/model"
	"acsp/internal/repository"
)

type ProjectModulesService struct {
	repo repository.ProjectModules
}

func NewProjectModulesService(repo repository.ProjectModules) *ProjectModulesService {
	return &ProjectModulesService{repo: repo}
}

func (p ProjectModulesService) Create(ctx context.Context, projectID int, module dto.CreateProjectModule) error {
	m := model.ProjectModule{
		ProjectID: projectID,
		Title:     module.Title,
	}

	err := p.repo.Create(ctx, m)
	if err != nil {
		return errors.Wrap(err, "Error occurred when creating project module")
	}

	return nil
}

func (p ProjectModulesService) Update(ctx context.Context, projectID, moduleID int, module dto.UpdateProjectModule) error {
	m := model.ProjectModule{
		ID:        moduleID,
		ProjectID: projectID,
		Title:     module.Title,
	}

	err := p.repo.Update(ctx, m)
	if err != nil {
		return errors.Wrap(err, "Error occurred when updating project module")
	}

	return nil
}

func (p ProjectModulesService) Delete(ctx context.Context, projectID, moduleID int) error {
	err := p.repo.Delete(ctx, projectID, moduleID)
	if err != nil {
		return errors.Wrap(err, "Error occurred when deleting project module")
	}

	return nil
}

func (p ProjectModulesService) GetAll(ctx context.Context, projectID int) ([]model.ProjectModule, error) {
	m, err := p.repo.GetAllByProjectID(ctx, projectID)
	if err != nil {
		return []model.ProjectModule{}, errors.Wrap(err, "Error occurred when getting project modules")
	}

	return m, nil
}

func (p ProjectModulesService) GetByID(ctx context.Context, moduleID int) (model.ProjectModule, error) {
	m, err := p.repo.GetByID(ctx, moduleID)
	if err != nil {
		return model.ProjectModule{}, errors.Wrap(err, "Error occurred when getting project module")
	}

	return m, nil
}
