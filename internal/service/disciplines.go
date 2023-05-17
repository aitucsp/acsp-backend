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

type DisciplinesService struct {
	repo         repository.Disciplines
	projectsRepo repository.Projects
}

func NewDisciplinesService(r repository.Disciplines, p repository.Projects) *DisciplinesService {
	return &DisciplinesService{repo: r, projectsRepo: p}
}

func (d DisciplinesService) Create(ctx context.Context, input dto.CreateDiscipline) error {
	l := logging.LoggerFromContext(ctx).With(zap.String(
		"disciplineTitle", input.Title),
	)

	discipline := model.Discipline{
		Title:       input.Title,
		Description: input.Description,
	}

	err := d.repo.Create(ctx, discipline)
	if err != nil {
		l.Error("Error when creating discipline", zap.Error(err))

		return errors.Wrap(err, "error when creating discipline")
	}

	return nil
}

func (d DisciplinesService) Update(ctx context.Context, disciplineID int, discipline dto.UpdateDiscipline) error {
	l := logging.LoggerFromContext(ctx).With(zap.String(
		"disciplineTitle", discipline.Title),
	)

	disciplineToUpdate := model.Discipline{
		ID:          disciplineID,
		Title:       discipline.Title,
		Description: discipline.Description,
	}

	err := d.repo.Update(ctx, disciplineToUpdate)
	if err != nil {
		l.Error("Error when updating discipline", zap.Error(err))

		return errors.Wrap(err, "error when updating discipline")
	}

	return nil
}

func (d DisciplinesService) Delete(ctx context.Context, disciplineID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int(
		"disciplineID", disciplineID),
	)

	err := d.repo.Delete(ctx, disciplineID)
	if err != nil {
		l.Error("Error when deleting discipline", zap.Error(err))

		return errors.Wrap(err, "error when deleting discipline")
	}

	return nil
}

func (d DisciplinesService) GetAll(ctx context.Context) ([]model.Discipline, error) {
	l := logging.LoggerFromContext(ctx)

	disciplines, err := d.repo.GetAll(ctx)
	if err != nil {
		l.Error("Error when getting all disciplines", zap.Error(err))

		return nil, errors.Wrap(err, "error when getting all disciplines")
	}

	return disciplines, nil
}

func (d DisciplinesService) GetByID(ctx context.Context, disciplineID int) (model.Discipline, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int(
		"disciplineID", disciplineID),
	)

	discipline, err := d.repo.GetByID(ctx, disciplineID)
	if err != nil {
		l.Error("Error when getting discipline by ID", zap.Error(err))

		return model.Discipline{}, errors.Wrap(err, "error when getting discipline by ID")
	}

	projects, err := d.projectsRepo.GetAllByDisciplineID(ctx, disciplineID)
	if err != nil {
		l.Error("Error when getting projects by discipline ID", zap.Error(err))

		return model.Discipline{}, errors.Wrap(err, "error when getting projects by discipline ID")
	}

	discipline.Projects = projects

	return discipline, nil
}
