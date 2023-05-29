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

type CourseModulesDatabase struct {
	db *sqlx.DB
}

func NewCourseModulesRepository(db *sqlx.DB) *CourseModulesDatabase {
	return &CourseModulesDatabase{
		db: db,
	}
}

func (c *CourseModulesDatabase) Create(ctx context.Context, module model.CourseModule) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("courseID", module.CourseID),
		zap.String("moduleTitle", module.Title),
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s (course_id, title, expected_result) 
								 VALUES ($1, $2, $3) RETURNING id`,
		constants.CourseModulesTable)

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

	res, err := stmt.Exec(module.CourseID, module.Title, module.ExpectedResult)
	if err != nil {
		l.Error("Error when executing the course module creating statement", zap.Error(err))

		return errors.Wrap(err, "error when executing the course module creating statement")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected of the course module", zap.Error(err))

		return errors.Wrap(err, "error when getting the affected of the course module")
	}

	if rowsAffected == 0 {
		l.Error("Error when getting the affected of the course module", zap.Error(apperror.ErrCreatingModule))

		return apperror.ErrCreatingModule
	}

	return nil
}

func (c *CourseModulesDatabase) Update(ctx context.Context, module model.CourseModule) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("moduleID", module.ID), zap.Int("courseID", module.CourseID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s SET title = $1, expected_result = $2, updated_at = now() 
											WHERE id = $3`,
		constants.CourseModulesTable)

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

	res, err := stmt.ExecContext(ctx, query, module.Title, module.ExpectedResult, module.ID)
	if err != nil {
		l.Error("Error when update the course module in database", zap.Error(err))

		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when updating the course module in database", zap.Error(err))

		return errors.Wrap(err, "error when updating the course module in database")
	}

	return nil
}

func (c *CourseModulesDatabase) Delete(ctx context.Context, courseID, moduleID int) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("courseID", courseID),
		zap.Int("moduleID", moduleID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND course_id = $2", constants.CourseModulesTable)

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

	res, err := stmt.ExecContext(ctx, query, moduleID, courseID)
	if err != nil {
		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when deleting the course module in database", zap.Error(err))

		return errors.Wrap(err, "error when deleting the course module in database")
	}

	return nil
}

func (c *CourseModulesDatabase) GetAllByCourseID(ctx context.Context, courseID int) ([]model.CourseModule, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("courseID", courseID))

	var modules []model.CourseModule

	query := fmt.Sprintf("SELECT * FROM %s WHERE project_id = $1",
		constants.CourseModulesTable)

	err := c.db.Select(&modules, query, courseID)
	if err != nil {
		l.Error("Error when getting the course's modules", zap.Error(err))

		return nil, errors.Wrap(err, "error when getting the course modules")
	}

	return modules, nil
}

func (c *CourseModulesDatabase) GetByID(ctx context.Context, moduleID int) (model.CourseModule, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("moduleID", moduleID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var module model.CourseModule

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		constants.CourseModulesTable)

	err := c.db.GetContext(ctx, &module, query, moduleID)
	if err != nil {
		l.Error("Error when getting the course module", zap.Error(err))

		return model.CourseModule{}, errors.Wrap(err, "error when executing the query")
	}

	return module, nil
}
