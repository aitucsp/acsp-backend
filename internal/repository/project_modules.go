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

type ProjectModulesDatabase struct {
	db *sqlx.DB
}

func NewProjectModulesRepository(db *sqlx.DB) *ProjectModulesDatabase {
	return &ProjectModulesDatabase{
		db: db,
	}
}

func (p *ProjectModulesDatabase) Create(ctx context.Context, module model.ProjectModule) error {
	l := logging.LoggerFromContext(ctx).With(
		zap.Int("projectID", module.ProjectID),
		zap.String("moduleTitle", module.Title),
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s (project_id, title) 
								 VALUES ($1, $2) RETURNING id`,
		constants.CodingLabProjectModulesTable)

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return err
	}

	res, err := stmt.Exec(module.ProjectID, module.Title)
	if err != nil {
		l.Error("Error when executing the project module creating statement", zap.Error(err))

		return err
	}

	id, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the id of the module", zap.Error(err))

		return err
	}

	if id == 0 {
		l.Error("Error when getting the id of the module", zap.Error(apperror.ErrCreatingCard))

		return apperror.ErrCreatingCard
	}

	return nil
}

func (p *ProjectModulesDatabase) Update(ctx context.Context, module model.ProjectModule) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("moduleID", module.ID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s SET title = $1, updated_at = now() 
											WHERE id = $2`,
		constants.CodingLabProjectsTable)

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return errors.Wrap(err, "error when preparing the query")
	}

	res, err := stmt.ExecContext(ctx, query, module.Title, module.ID)
	if err != nil {
		l.Error("Error when update the project module in database", zap.Error(err))

		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when updating the project module in database", zap.Error(err))

		return errors.Wrap(err, "error when updating the project module in database")
	}

	return nil
}

func (p *ProjectModulesDatabase) Delete(ctx context.Context, projectID, moduleID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("moduleID", moduleID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND project_id = $2", constants.CodingLabProjectModulesTable)

	stmt, err := p.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "error when preparing the query")
	}

	res, err := stmt.ExecContext(ctx, query, moduleID, projectID)
	if err != nil {
		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when deleting the project module in database", zap.Error(err))

		return errors.Wrap(err, "error when deleting the project module in database")
	}

	return nil
}

func (p *ProjectModulesDatabase) GetAllByProjectID(ctx context.Context, projectID int) ([]model.ProjectModule, error) {
	var modules []model.ProjectModule

	query := fmt.Sprintf("SELECT * FROM %s WHERE project_id = $1",
		constants.CodingLabProjectModulesTable)

	err := p.db.Select(&modules, query, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "error when getting the projects")
	}

	return modules, nil
}

func (p *ProjectModulesDatabase) GetByID(ctx context.Context, moduleID int) (model.ProjectModule, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("moduleID", moduleID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var module model.ProjectModule

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		constants.CodingLabProjectsTable)

	err := p.db.GetContext(ctx, &module, query, moduleID)
	if err != nil {
		l.Error("Error when getting the project module", zap.Error(err))

		return model.ProjectModule{}, errors.Wrap(err, "error when executing the query")
	}

	return module, nil
}
