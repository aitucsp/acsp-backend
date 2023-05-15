package service

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"acsp/internal/dto"
	"acsp/internal/logging"
	"acsp/internal/model"
	"acsp/internal/repository"
)

type ProjectModulesService struct {
	repo repository.ProjectModules
}

func NewProjectModulesService(repo repository.ProjectModules) *ProjectModulesService {
	return &ProjectModulesService{repo: repo}
}

func (p ProjectModulesService) Create(ctx context.Context, projectID int, input dto.CreateProjectModule) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.String("projectModuleTitle", input.Title),
		zap.Int("projectID", projectID),
	)

	m := model.ProjectModule{
		ProjectID:    projectID,
		Title:        input.Title,
		Description:  input.Description,
		ReferenceURL: input.ReferenceURL,
	}

	err := p.repo.Create(ctx, m)
	if err != nil {
		l.Error("Error when creating a project module", zap.Error(err))

		return errors.Wrap(err, "Error occurred when creating project module")
	}

	return nil
}

func (p ProjectModulesService) Update(ctx context.Context, projectID, moduleID int, input dto.UpdateProjectModule) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.String("projectModuleTitle", input.Title),
		zap.Int("projectID", projectID),
		zap.Int("moduleID", moduleID),
	)

	m := model.ProjectModule{
		ID:           moduleID,
		ProjectID:    projectID,
		Title:        input.Title,
		Description:  input.Description,
		ReferenceURL: input.ReferenceURL,
	}

	err := p.repo.Update(ctx, m)
	if err != nil {
		l.Error("Error when updating a project module", zap.Error(err))

		return errors.Wrap(err, "Error occurred when updating project module")
	}

	return nil
}

func (p ProjectModulesService) Delete(ctx context.Context, projectID, moduleID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("projectID", projectID),
		zap.Int("moduleID", moduleID),
	)

	err := p.repo.Delete(ctx, projectID, moduleID)
	if err != nil {
		l.Error("Error when deleting a project module", zap.Error(err))

		return errors.Wrap(err, "Error occurred when deleting project module")
	}

	return nil
}

func (p ProjectModulesService) GetAll(ctx context.Context, projectID int) ([]model.ProjectModule, error) {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("projectID", projectID),
	)

	m, err := p.repo.GetAllByProjectID(ctx, projectID)
	if err != nil {
		l.Error("Error when getting project modules", zap.Error(err))

		return []model.ProjectModule{}, errors.Wrap(err, "Error occurred when getting project modules")
	}

	return m, nil
}

func (p ProjectModulesService) GetByID(ctx context.Context, moduleID int) (model.ProjectModule, error) {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("moduleID", moduleID),
	)

	m, err := p.repo.GetByID(ctx, moduleID)
	if err != nil {
		l.Error("Error when getting project module", zap.Error(err))

		return model.ProjectModule{}, errors.Wrap(err, "Error occurred when getting project module")
	}

	return m, nil
}
