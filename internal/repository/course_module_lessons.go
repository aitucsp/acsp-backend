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

type CourseModuleLessonsDatabase struct {
	db *sqlx.DB
}

func NewCourseModuleLessonsRepository(db *sqlx.DB) *CourseModuleLessonsDatabase {
	return &CourseModuleLessonsDatabase{
		db: db,
	}
}

func (c *CourseModuleLessonsDatabase) Create(ctx context.Context, lesson model.CourseModuleLesson) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("moduleID", lesson.ModuleID),
		zap.String("lessonTitle", lesson.Title),
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s (module_id, title, description, reference_url) 
								 VALUES ($1, $2, $3, $4) RETURNING id`,
		constants.CourseModuleLessonsTable)

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

	res, err := stmt.Exec(lesson.ModuleID, lesson.Title, lesson.Description, lesson.ReferenceURL)
	if err != nil {
		l.Error("Error when executing the module lesson creating statement", zap.Error(err))

		return errors.Wrap(err, "error when executing the module lesson creating statement")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected of the module lesson", zap.Error(err))

		return errors.Wrap(err, "error when getting the affected of the module lesson")
	}

	if rowsAffected == 0 {
		l.Error("Error when getting the affected of the module lesson", zap.Error(apperror.ErrCreatingModule))

		return apperror.ErrCreatingModule
	}

	return nil
}

func (c *CourseModuleLessonsDatabase) Update(ctx context.Context, lesson model.CourseModuleLesson) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("lessonID", lesson.ID),
		zap.Int("courseID", lesson.ModuleID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s SET title = $1, description = $2, reference_url = $3,
								 updated_at = NOW() WHERE id = $4`,
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

	res, err := stmt.ExecContext(ctx, query, lesson.Title, lesson.Description, lesson.ReferenceURL, lesson.ID)
	if err != nil {
		l.Error("Error when update the module lesson in database", zap.Error(err))

		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when updating the module lesson in database", zap.Error(err))

		return errors.Wrap(err, "error when updating the module lesson in database")
	}

	return nil
}

func (c *CourseModuleLessonsDatabase) Delete(ctx context.Context, lessonID, moduleID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("courseID", lessonID),
		zap.Int("moduleID", moduleID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND module_Id = $2", constants.CourseModuleLessonsTable)

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

	res, err := stmt.ExecContext(ctx, query, lessonID, moduleID)
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

func (c *CourseModuleLessonsDatabase) GetAllByModuleId(ctx context.Context, moduleID int) ([]model.CourseModuleLesson, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("moduleID", moduleID))

	var modules []model.CourseModuleLesson

	query := fmt.Sprintf("SELECT * FROM %s WHERE module_id = $1",
		constants.CourseModuleLessonsTable)

	err := c.db.Select(&modules, query, moduleID)
	if err != nil {
		l.Error("Error when getting the module's lessons", zap.Error(err))

		return nil, errors.Wrap(err, "error when getting the module lessons")
	}

	return modules, nil
}

func (c *CourseModuleLessonsDatabase) GetByID(ctx context.Context, lessonID int) (model.CourseModuleLesson, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("lessonID", lessonID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var lesson model.CourseModuleLesson

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		constants.CourseModuleLessonsTable)

	err := c.db.GetContext(ctx, &lesson, query, lessonID)
	if err != nil {
		l.Error("Error when getting the module lesson", zap.Error(err))

		return model.CourseModuleLesson{}, errors.Wrap(err, "error when executing the query")
	}

	return lesson, nil
}
