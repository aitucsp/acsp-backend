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

type ProjectsDatabase struct {
	db *sqlx.DB
}

func NewProjectsRepository(db *sqlx.DB) *ProjectsDatabase {
	return &ProjectsDatabase{
		db: db,
	}
}

func (p *ProjectsDatabase) Create(ctx context.Context, disciplineID int, project model.Project) error {
	l := logging.LoggerFromContext(ctx).With(zap.String("projectTitle", project.Title))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s 
								(discipline_id, title, description, level, image_url, work_hours)
								 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		constants.CodingLabProjectsTable)

	stmt, err := p.db.PrepareContext(ctx, query)
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
		disciplineID,
		project.Title,
		project.Description,
		project.Level,
		project.ImageURL,
		project.WorkHours,
	)
	if err != nil {
		l.Error("Error when executing the project creating statement", zap.Error(err))

		return err
	}

	id, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the id of the project", zap.Error(err))

		return err
	}

	if id == 0 {
		l.Error("Error when getting the id of the project", zap.Error(apperror.ErrCreatingCard))

		return apperror.ErrCreatingCard
	}

	return nil
}

func (p *ProjectsDatabase) Update(ctx context.Context, project model.Project) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("projectID", project.ID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s SET title = $1, description = $2, level = $3, work_hours = $4, updated_at = now() 
											WHERE id = $5`,
		constants.CodingLabProjectsTable)

	stmt, err := p.db.PrepareContext(ctx, query)
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
		project.Title,
		project.Description,
		project.Level,
		project.WorkHours,
		project.ID,
	)
	if err != nil {
		l.Error("Error when update the project in database", zap.Error(err))

		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when updating the project in database", zap.Error(err))

		return errors.Wrap(err, "error when updating the project in database")
	}

	return nil
}

func (p *ProjectsDatabase) Delete(ctx context.Context, projectID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("projectID", projectID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", constants.CodingLabProjectsTable)

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "error when preparing the query")
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.ExecContext(ctx, query, projectID)
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

		return errors.Wrap(err, "error when deleting the project in database")
	}

	return nil
}

func (p *ProjectsDatabase) GetAll(ctx context.Context) ([]model.Project, error) {
	var projects []model.Project

	query := fmt.Sprintf("SELECT * FROM %s",
		constants.CodingLabProjectsTable)

	err := p.db.Select(&projects, query)
	if err != nil {
		return nil, errors.Wrap(err, "error when getting the projects")
	}

	return projects, nil
}

func (p *ProjectsDatabase) GetByID(ctx context.Context, projectID int) (model.Project, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("projectID", projectID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var project model.Project

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		constants.CodingLabProjectsTable)

	err := p.db.GetContext(ctx, &project, query, projectID)
	if err != nil {
		l.Error("Error when getting the project", zap.Error(err))

		return model.Project{}, errors.Wrap(err, "error when executing the query")
	}

	return project, nil
}
