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

type CourseModuleLessonsService struct {
	repo repository.CourseLessons
}

func NewCourseModuleLessonsService(repo repository.CourseLessons) *CourseModuleLessonsService {
	return &CourseModuleLessonsService{repo: repo}
}

func (c CourseModuleLessonsService) Create(ctx context.Context, moduleID int, input dto.CreateCourseModuleLesson) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.String("courseModuleLessonTitle", input.Title),
		zap.Int("moduleID", moduleID),
	)

	m := model.CourseModuleLesson{
		ModuleID:     moduleID,
		Title:        input.Title,
		Description:  input.Description,
		ReferenceURL: input.ReferenceURL,
	}

	err := c.repo.Create(ctx, m)
	if err != nil {
		l.Error("Error when creating a course module lesson", zap.Error(err))

		return errors.Wrap(err, "Error occurred when creating course module lesson")
	}

	return nil
}

func (c CourseModuleLessonsService) Update(ctx context.Context, moduleID, lessonID int, lesson dto.UpdateCourseModuleLesson) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.String("courseModuleLessonTitle", lesson.Title),
		zap.Int("moduleID", moduleID),
		zap.Int("lessonID", lessonID),
	)

	m := model.CourseModuleLesson{
		ID:           lessonID,
		ModuleID:     moduleID,
		Title:        lesson.Title,
		Description:  lesson.Description,
		ReferenceURL: lesson.ReferenceURL,
	}

	err := c.repo.Update(ctx, m)
	if err != nil {
		l.Error("Error when updating a course module lesson", zap.Error(err))

		return errors.Wrap(err, "Error occurred when updating course module lesson")
	}

	return nil
}

func (c CourseModuleLessonsService) Delete(ctx context.Context, moduleID, lessonID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("moduleID", moduleID),
		zap.Int("lessonID", lessonID),
	)

	err := c.repo.Delete(ctx, moduleID, lessonID)
	if err != nil {
		l.Error("Error when deleting a course module lesson", zap.Error(err))

		return errors.Wrap(err, "Error occurred when deleting course module lesson")
	}

	return nil
}

func (c CourseModuleLessonsService) GetAll(ctx context.Context, moduleID int) ([]model.CourseModuleLesson, error) {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("moduleID", moduleID),
	)

	lessons, err := c.repo.GetAllByModuleId(ctx, moduleID)
	if err != nil {
		l.Error("Error when getting all course module lessons", zap.Error(err))

		return nil, errors.Wrap(err, "Error occurred when getting all course module lessons")
	}

	return lessons, nil
}

func (c CourseModuleLessonsService) GetByID(ctx context.Context, lessonID int) (model.CourseModuleLesson, error) {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("lessonID", lessonID),
	)

	lesson, err := c.repo.GetByID(ctx, lessonID)
	if err != nil {
		l.Error("Error when getting course module lesson by id", zap.Error(err))

		return model.CourseModuleLesson{}, errors.Wrap(err, "Error occurred when getting course module lesson by id")
	}

	return lesson, nil
}
