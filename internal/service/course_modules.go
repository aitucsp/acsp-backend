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

type CourseModulesService struct {
	repo repository.CourseModules
}

func NewCourseModulesService(repo repository.CourseModules) *CourseModulesService {
	return &CourseModulesService{repo: repo}
}

func (c *CourseModulesService) Create(ctx context.Context, courseID int, input dto.CreateCourseModule) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.String("courseModuleTitle", input.Title),
		zap.Int("courseID", courseID),
	)

	m := model.CourseModule{
		CourseID:       courseID,
		Title:          input.Title,
		ExpectedResult: input.ExpectedResult,
	}

	err := c.repo.Create(ctx, m)
	if err != nil {
		l.Error("Error when creating a project module", zap.Error(err))

		return errors.Wrap(err, "Error occurred when creating project module")
	}

	return nil
}

func (c *CourseModulesService) Update(ctx context.Context, courseID, moduleID int, input dto.UpdateCourseModule) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.String("projectModuleTitle", input.Title),
		zap.Int("courseID", courseID),
		zap.Int("moduleID", moduleID),
	)

	m := model.CourseModule{
		ID:             moduleID,
		CourseID:       courseID,
		Title:          input.Title,
		ExpectedResult: input.ExpectedResult,
	}

	err := c.repo.Update(ctx, m)
	if err != nil {
		l.Error("Error when updating a course module", zap.Error(err))

		return errors.Wrap(err, "Error occurred when updating course module")
	}

	return nil
}

func (c *CourseModulesService) Delete(ctx context.Context, courseID, moduleID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("courseID", courseID),
		zap.Int("moduleID", moduleID),
	)

	err := c.repo.Delete(ctx, courseID, moduleID)
	if err != nil {
		l.Error("Error when deleting a course module", zap.Error(err))

		return errors.Wrap(err, "Error occurred when deleting course module")
	}

	return nil
}

func (c *CourseModulesService) GetAll(ctx context.Context, courseID int) ([]model.CourseModule, error) {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("courseID", courseID),
	)

	modules, err := c.repo.GetAllByCourseID(ctx, courseID)
	if err != nil {
		l.Error("Error when getting all course modules", zap.Error(err))

		return nil, errors.Wrap(err, "Error occurred when getting all course modules")
	}

	return modules, nil
}

func (c *CourseModulesService) GetByID(ctx context.Context, moduleID int) (model.CourseModule, error) {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("moduleID", moduleID),
	)

	module, err := c.repo.GetByID(ctx, moduleID)
	if err != nil {
		l.Error("Error when getting course module by ID", zap.Error(err))

		return model.CourseModule{}, errors.Wrap(err, "Error occurred when getting course module by ID")
	}

	return module, nil
}
