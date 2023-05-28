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

type ContestsDatabase struct {
	db *sqlx.DB
}

func NewContestsRepository(db *sqlx.DB) *ContestsDatabase {
	return &ContestsDatabase{
		db: db,
	}
}

func (c *ContestsDatabase) Create(ctx context.Context, contest model.Contest) error {
	l := logging.LoggerFromContext(ctx).With(zap.String("contestName", contest.Name))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s 
								(contest_name, description, link, start_date, end_date)
								 VALUES ($1, $2, $3, $4) RETURNING rowsAffected`,
		constants.ContestsTable)

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
		contest.Name,
		contest.Description,
		contest.Link,
		contest.StartDate,
		contest.EndDate,
	)
	if err != nil {
		l.Error("Error when executing the contest creating statement", zap.Error(err))

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rowsAffected of the creating contest query", zap.Error(err))

		return err
	}

	if rowsAffected == 0 {
		l.Error("Error when getting the rows affected of the creating contest query",
			zap.Error(apperror.ErrCreatingContest))

		return apperror.ErrCreatingContest
	}

	return nil
}

func (c *ContestsDatabase) Update(ctx context.Context, contest model.Contest) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("contestID", contest.ID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s SET contest_name = $1, 
												description = $2, 
												link = $3,
												start_date = $4
												end_date = $5, 
												updated_at = now() 
											WHERE id = $6`,
		constants.CodingLabProjectsTable)

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
		contest.Name,
		contest.Description,
		contest.Link,
		contest.StartDate,
		contest.EndDate,
		contest.ID,
	)
	if err != nil {
		l.Error("Error when update the contest in database", zap.Error(err))

		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when updating the contest in database", zap.Error(err))

		return errors.Wrap(err, "error when updating the contest in database")
	}

	return nil
}

func (c *ContestsDatabase) Delete(ctx context.Context, contestID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("contestID", contestID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", constants.ContestsTable)

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

	res, err := stmt.ExecContext(ctx, query, contestID)
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

func (c *ContestsDatabase) GetByID(ctx context.Context, contestID int) (model.Contest, error) {
	var contest model.Contest

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		constants.ContestsTable)

	err := c.db.Get(&contest, query, contestID)
	if err != nil {
		return model.Contest{}, errors.Wrap(err, "error when executing the query")
	}

	return contest, nil
}

func (c *ContestsDatabase) GetAll(ctx context.Context) ([]model.Contest, error) {
	var contests []model.Contest

	query := fmt.Sprintf("SELECT * FROM %s",
		constants.ContestsTable)

	err := c.db.Select(&contests, query)
	if err != nil {
		return nil, errors.Wrap(err, "error when getting the contest")
	}

	return contests, nil
}
