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

type DisciplinesDatabase struct {
	db *sqlx.DB
}

func NewDisciplinesRepository(db *sqlx.DB) *DisciplinesDatabase {
	return &DisciplinesDatabase{
		db: db,
	}
}

func (d *DisciplinesDatabase) Create(ctx context.Context, discipline model.Discipline) error {
	l := logging.LoggerFromContext(ctx).With(zap.String(
		"disciplineTitle", discipline.Title),
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s 
								(title, description)
								 VALUES ($1, $2)`,
		constants.CodingLabDisciplinesTable)

	stmt, err := d.db.PrepareContext(ctx, query)
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

	res, err := stmt.Exec(discipline.Title, discipline.Description)
	if err != nil {
		l.Error("Error when executing the discipline creating statement", zap.Error(err))

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rowsAffected of the discipline", zap.Error(err))

		return err
	}

	if rowsAffected == 0 {
		l.Error("Error when getting the rowsAffected of the project", zap.Error(apperror.ErrCreatingDiscipline))

		return apperror.ErrCreatingDiscipline
	}

	return nil
}

func (d *DisciplinesDatabase) Update(ctx context.Context, discipline model.Discipline) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("disciplineID", discipline.ID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s SET title = $1, description = $2, image_url = $3, updated_at = now() 
											WHERE id = $4`,
		constants.CodingLabDisciplinesTable)

	stmt, err := d.db.PrepareContext(ctx, query)
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
		discipline.Title,
		discipline.Description,
		discipline.ImageURL,
		discipline.ID,
	)
	if err != nil {
		l.Error("Error when update the discipline in database", zap.Error(err))

		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when updating the discipline in database", zap.Error(err))

		return errors.Wrap(err, "error when updating the discipline in database")
	}

	return nil
}

func (d *DisciplinesDatabase) Delete(ctx context.Context, disciplineID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("disciplineID", disciplineID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", constants.CodingLabDisciplinesTable)

	stmt, err := d.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "error when preparing the query")
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.ExecContext(ctx, query, disciplineID)
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

func (d *DisciplinesDatabase) GetByID(ctx context.Context, disciplineID int) (model.Discipline, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("disciplineID", disciplineID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", constants.CodingLabDisciplinesTable)

	stmt, err := d.db.PrepareContext(ctx, query)
	if err != nil {
		l.Error("Error when preparing the query", zap.Error(err))

		return model.Discipline{}, errors.Wrap(err, "error when preparing the query")
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	var discipline model.Discipline
	err = stmt.QueryRowContext(ctx, query, disciplineID).Scan(
		&discipline.ID,
		&discipline.Title,
		&discipline.Description,
		&discipline.ImageURL,
		&discipline.CreatedAt,
		&discipline.UpdatedAt,
	)
	if err != nil {
		l.Error("Error when scanning the discipline", zap.Error(err))

		return model.Discipline{}, errors.Wrap(err, "error when scanning the discipline")
	}

	return discipline, nil
}

func (d *DisciplinesDatabase) GetAll(ctx context.Context) ([]model.Discipline, error) {
	l := logging.LoggerFromContext(ctx)

	var disciplines []model.Discipline

	query := fmt.Sprintf("SELECT * FROM %s",
		constants.CodingLabDisciplinesTable)

	err := d.db.Select(&disciplines, query)
	if err != nil {
		l.Error("Error when getting the disciplines", zap.Error(err))

		return nil, errors.Wrap(err, "error when getting the disciplines")
	}

	return disciplines, nil
}
