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

type CoursesService struct {
	repo        repository.Courses
	modulesRepo repository.CourseModules
}

func NewCoursesService(repo repository.Courses, modulesRepo repository.CourseModules) *CoursesService {
	return &CoursesService{repo: repo, modulesRepo: modulesRepo}
}

func (c *CoursesService) Create(ctx context.Context, input dto.CreateCourse) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.String("courseTitle", input.Title),
	)

	course := model.Course{
		Title:       input.Title,
		Description: input.Description,
	}

	err := c.repo.Create(ctx, course)
	if err != nil {
		l.Error("Error when creating a course", zap.Error(err))

		return errors.Wrap(err, "error when creating a course")
	}

	return nil
}

func (c *CoursesService) Update(ctx context.Context, courseID int, input dto.UpdateCourse) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("courseID", courseID),
	)

	project := model.Course{
		ID:          courseID,
		Title:       input.Title,
		Description: input.Description,
		Rating:      input.Rating,
	}

	err := c.repo.Update(ctx, project)
	if err != nil {
		l.Error("Error when updating a course", zap.Error(err))

		return errors.Wrap(err, "error when updating a course")
	}

	return nil
}

func (c *CoursesService) Delete(ctx context.Context, courseID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("courseID", courseID),
	)

	err := c.repo.Delete(ctx, courseID)
	if err != nil {
		l.Error("Error when deleting a course", zap.Error(err))

		return errors.Wrap(err, "error when deleting a course")
	}

	return nil
}

func (c *CoursesService) GetAll(ctx context.Context) ([]model.Course, error) {
	l := logging.LoggerFromContext(ctx)

	courses, err := c.repo.GetAll(ctx)
	if err != nil {
		l.Error("Error when getting all courses", zap.Error(err))

		return nil, errors.Wrap(err, "error when getting all courses")
	}

	return courses, nil
}

func (c *CoursesService) GetByID(ctx context.Context, courseID int) (model.Course, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("courseID", courseID))

	course, err := c.repo.GetByID(ctx, courseID)
	if err != nil {
		l.Error("Error when getting a course by ID", zap.Error(err))

		return model.Course{}, errors.Wrap(err, "error when getting a project by ID")
	}

	m, err := c.modulesRepo.GetAllByCourseID(ctx, courseID)
	if err != nil {
		l.Error("Error when getting all modules by course ID", zap.Error(err))

		return model.Course{}, errors.Wrap(err, "error when getting all modules by course ID")
	}

	if m == nil {
		m = []model.CourseModule{}
	} else {
		course.Modules = m
	}

	return course, nil
}
