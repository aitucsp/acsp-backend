package repository

import (
	"context"
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

func (p *ProjectsDatabase) Create(ctx context.Context, project model.Project) error {
	l := logging.LoggerFromContext(ctx).With(zap.String("projectTitle", project.Title))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s (title, description, reference_list) 
								 VALUES ($1, $2, $3) RETURNING id`,
		constants.CodingLabProjectsTable)

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return err
	}

	res, err := stmt.Exec(project.Title, project.Description, project.ReferenceList)
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

	query := fmt.Sprintf(`UPDATE %s SET title = $1, description = $2, reference_list = $3 updated_at = now() 
											WHERE id = $4`,
		constants.CodingLabProjectsTable)

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return errors.Wrap(err, "error when preparing the query")
	}

	res, err := stmt.ExecContext(ctx, query, project.Title, project.Description, project.ReferenceList, project.ID)
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
