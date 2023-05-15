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

type ProjectsService struct {
	repo        repository.Projects
	modulesRepo repository.ProjectModules
}

func NewProjectsService(repo repository.Projects, modulesRepo repository.ProjectModules) *ProjectsService {
	return &ProjectsService{repo: repo, modulesRepo: modulesRepo}
}

func (p *ProjectsService) Create(ctx context.Context, disciplineID int, input dto.CreateProject) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.String("projectTitle", input.Title),
		zap.Int("disciplineID", disciplineID),
	)

	project := model.Project{
		DisciplineID: disciplineID,
		Title:        input.Title,
		Description:  input.Description,
		Level:        input.Level,
		WorkHours:    input.WorkHours,
	}

	err := p.repo.Create(ctx, disciplineID, project)
	if err != nil {
		l.Error("Error when creating a project", zap.Error(err))

		return errors.Wrap(err, "error when creating a project")
	}

	return nil
}

func (p *ProjectsService) Update(ctx context.Context, input dto.UpdateProject) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.String("projectTitle", input.Title),
	)

	project := model.Project{
		Title:       input.Title,
		Description: input.Description,
		Level:       input.Level,
		WorkHours:   input.WorkHours,
	}

	err := p.repo.Update(ctx, project)
	if err != nil {
		l.Error("Error when updating a project", zap.Error(err))

		return errors.Wrap(err, "error when updating a project")
	}

	return nil
}

func (p *ProjectsService) Delete(ctx context.Context, projectID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("projectID", projectID))

	err := p.repo.Delete(ctx, projectID)
	if err != nil {
		l.Error("Error when deleting a project", zap.Error(err))

		return errors.Wrap(err, "error when deleting a project")
	}

	return nil
}

func (p *ProjectsService) GetAll(ctx context.Context) ([]model.Project, error) {
	l := logging.LoggerFromContext(ctx)

	projects, err := p.repo.GetAll(ctx)
	if err != nil {
		l.Error("Error when getting all projects", zap.Error(err))

		return nil, errors.Wrap(err, "error when getting all projects")
	}

	return projects, nil
}

func (p *ProjectsService) GetByID(ctx context.Context, projectID int) (model.Project, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("projectID", projectID))

	project, err := p.repo.GetByID(ctx, projectID)
	if err != nil {
		l.Error("Error when getting a project by ID", zap.Error(err))

		return model.Project{}, errors.Wrap(err, "error when getting a project by ID")
	}

	m, err := p.modulesRepo.GetAllByProjectID(ctx, projectID)
	if err != nil {
		l.Error("Error when getting all modules by project ID", zap.Error(err))

		return model.Project{}, errors.Wrap(err, "error when getting all modules by project ID")
	}

	if m == nil {
		m = []model.ProjectModule{}
	} else {
		project.Modules = m
	}
	return project, nil
}
