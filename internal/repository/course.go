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

type CoursesDatabase struct {
	db *sqlx.DB
}

func NewCoursesRepository(db *sqlx.DB) *CoursesDatabase {
	return &CoursesDatabase{
		db: db,
	}
}

func (c *CoursesDatabase) Create(ctx context.Context, course model.Course) error {
	l := logging.LoggerFromContext(ctx).With(zap.String("courseTitle", course.Title))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s 
								(author_id, title, description, rating)
								 VALUES ($1, $2, $3, $4) RETURNING id`,
		constants.CoursesTable)

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

	res, err := stmt.Exec(
		course.AuthorID,
		course.Title,
		course.Description,
		course.Rating,
	)
	if err != nil {
		l.Error("Error when executing the course creating statement", zap.Error(err))

		return errors.Wrap(err, "error when executing the course creating statement")
	}

	id, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the id of the course", zap.Error(err))

		return errors.Wrap(err, "error when getting the id of the course")
	}

	if id == 0 {
		l.Error("Error when getting the id of the course", zap.Error(apperror.ErrCreatingCourse))

		return apperror.ErrCreatingCourse
	}

	return nil
}

func (c *CoursesDatabase) Update(ctx context.Context, course model.Course) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("courseID", course.ID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s SET title = $1, description = $2, rating = $3, updated_at = now() 
											WHERE id = $4`,
		constants.CoursesTable)

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

	res, err := stmt.ExecContext(ctx, query,
		course.Title,
		course.Description,
		course.Rating,
		course.ID,
	)
	if err != nil {
		l.Error("Error when update the course in database", zap.Error(err))

		return errors.Wrap(err, "error when executing the query to update the course")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when updating the project in database", zap.Error(err))

		return errors.Wrap(err, "error when updating the course in database")
	}

	return nil
}

func (c *CoursesDatabase) Delete(ctx context.Context, courseID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("courseID", courseID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", constants.CoursesTable)

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

	res, err := stmt.ExecContext(ctx, query, courseID)
	if err != nil {
		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when deleting the project in database", zap.Error(err))

		return errors.Wrap(err, "error when deleting the course in database")
	}

	return nil
}

func (c *CoursesDatabase) GetAll(ctx context.Context) ([]model.Course, error) {
	var courses []model.Course

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s",
		constants.CoursesTable)

	err := c.db.Select(&courses, query)
	if err != nil {
		return nil, errors.Wrap(err, "error when getting the courses")
	}

	return courses, nil
}

func (c *CoursesDatabase) GetByID(ctx context.Context, courseID int) (model.Course, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("courseID", courseID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var course model.Course

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		constants.CodingLabProjectsTable)

	err := c.db.GetContext(ctx, &course, query, courseID)
	if err != nil {
		l.Error("Error when getting the course", zap.Error(err))

		return model.Course{}, errors.Wrap(err, "error when executing the query to get the course")
	}

	return course, nil
}
