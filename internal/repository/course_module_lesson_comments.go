package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"acsp/internal/apperror"
	"acsp/internal/constants"
	"acsp/internal/logging"
	"acsp/internal/model"
)

type CourseModuleLessonCommentsDatabase struct {
	db *sqlx.DB
}

func NewCourseModuleLessonCommentsRepository(db *sqlx.DB) *CourseModuleLessonCommentsDatabase {
	return &CourseModuleLessonCommentsDatabase{
		db: db,
	}
}

func (c *CourseModuleLessonCommentsDatabase) Create(ctx context.Context, comment model.CourseModuleLessonComment) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("lessonID", comment.LessonID),
		zap.String("text", comment.Text),
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s (lesson_id, user_id, text) 
								 VALUES ($1, $2, $3) RETURNING id`,
		constants.CourseLessonCommentsTable)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.Exec(comment.LessonID, comment.AuthorID, comment.Text)
	if err != nil {
		l.Error("Error when executing the lesson comment creating statement", zap.Error(err))

		return errors.Wrap(err, "error when executing the lesson comment creating statement")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected of the lesson comment", zap.Error(err))

		return errors.Wrap(err, "error when getting the affected of the lesson comment")
	}

	if rowsAffected == 0 {
		l.Error("Error when getting the affected of the lesson comment", zap.Error(apperror.ErrCreatingModule))

		return apperror.ErrCreatingModule
	}

	return nil
}

func (c *CourseModuleLessonCommentsDatabase) Update(ctx context.Context, comment model.CourseModuleLessonComment) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("commentID", comment.ID),
		zap.Int("lessonID", comment.LessonID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s SET text = $1, updated_at = NOW() WHERE id = $2`,
		constants.CourseModuleLessonsTable)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return errors.Wrap(err, "error when preparing the query")
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.ExecContext(ctx, query, comment.Text, comment.ID)
	if err != nil {
		l.Error("Error when update the lesson comment in database", zap.Error(err))

		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when updating the lesson comment in database", zap.Error(err))

		return errors.Wrap(err, "error when updating the lesson comment in database")
	}

	return nil
}

func (c *CourseModuleLessonCommentsDatabase) Delete(ctx context.Context, lessonID int, commentID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("lessonID", lessonID),
		zap.Int("commentID", commentID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND lesson_id = $2", constants.CourseLessonCommentsTable)

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "error when preparing the query")
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.ExecContext(ctx, query, commentID, lessonID)
	if err != nil {
		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when deleting the module lesson in database", zap.Error(err))

		return errors.Wrap(err, "error when deleting the lesson module in database")
	}

	return nil
}

func (c *CourseModuleLessonCommentsDatabase) GetAllByLessonId(ctx context.Context, lessonID int) ([]model.CourseModuleLessonComment, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("lessonID", lessonID))

	var comments []model.CourseModuleLessonComment

	query := fmt.Sprintf("SELECT * FROM %s WHERE lesson_id = $1",
		constants.CourseLessonCommentsTable)

	err := c.db.Select(&comments, query, lessonID)
	if err != nil {
		l.Error("Error when getting the lesson's comments", zap.Error(err))

		return nil, errors.Wrap(err, "error when getting the lesson comments")
	}

	return comments, nil
}

func (c *CourseModuleLessonCommentsDatabase) GetByID(ctx context.Context, commentID int) (model.CourseModuleLessonComment, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("commentID", commentID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var comment model.CourseModuleLessonComment

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		constants.CourseModuleLessonsTable)

	err := c.db.GetContext(ctx, &comment, query, commentID)
	if err != nil {
		l.Error("Error when getting the comment", zap.Error(err))

		return model.CourseModuleLessonComment{}, errors.Wrap(err, "error when executing the query")
	}

	return comment, nil
}
