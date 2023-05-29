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

type CourseModuleLessonCommentsService struct {
	repo repository.CourseLessonComments
}

func NewCourseModuleLessonCommentsService(repo repository.CourseLessonComments) *CourseModuleLessonCommentsService {
	return &CourseModuleLessonCommentsService{repo: repo}
}

func (c *CourseModuleLessonCommentsService) Create(ctx context.Context, lessonID int, input dto.CreateLessonComment) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("lessonID", lessonID),
		zap.String("comment", input.Text),
	)

	m := model.CourseModuleLessonComment{
		LessonID: lessonID,
		Text:     input.Text,
	}

	err := c.repo.Create(ctx, m)
	if err != nil {
		l.Error("Error when creating a course module lesson comment", zap.Error(err))

		return errors.Wrap(err, "Error occurred when creating course module lesson comment")
	}

	return nil
}

func (c *CourseModuleLessonCommentsService) Update(ctx context.Context, lessonID, commentID int, comment dto.UpdateLessonComment) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("lessonID", lessonID),
		zap.Int("commentID", commentID),
		zap.String("comment", comment.Text),
	)

	m := model.CourseModuleLessonComment{
		ID:       commentID,
		LessonID: lessonID,
		Text:     comment.Text,
	}

	err := c.repo.Update(ctx, m)
	if err != nil {
		l.Error("Error when updating a course module lesson comment", zap.Error(err))

		return errors.Wrap(err, "Error occurred when updating course module lesson comment")
	}

	return nil
}

func (c *CourseModuleLessonCommentsService) Delete(ctx context.Context, lessonID, commentID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("lessonID", lessonID),
		zap.Int("commentID", commentID),
	)

	err := c.repo.Delete(ctx, lessonID, commentID)
	if err != nil {
		l.Error("Error when deleting a course module lesson comment", zap.Error(err))

		return errors.Wrap(err, "Error occurred when deleting course module lesson comment")
	}

	return nil
}

func (c *CourseModuleLessonCommentsService) GetAll(ctx context.Context, lessonID int) ([]model.CourseModuleLessonComment, error) {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("lessonID", lessonID),
	)

	comments, err := c.repo.GetAllByLessonId(ctx, lessonID)
	if err != nil {
		l.Error("Error when getting all course module lesson comments", zap.Error(err))

		return nil, errors.Wrap(err, "Error occurred when getting all course module lesson comments")
	}

	return comments, nil
}

func (c *CourseModuleLessonCommentsService) GetByID(ctx context.Context, commentID int) (model.CourseModuleLessonComment, error) {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("commentID", commentID),
	)

	comment, err := c.repo.GetByID(ctx, commentID)
	if err != nil {
		l.Error("Error when getting a course module lesson comment", zap.Error(err))

		return model.CourseModuleLessonComment{}, errors.Wrap(err, "Error occurred when getting course module lesson comment")
	}

	return comment, nil
}
