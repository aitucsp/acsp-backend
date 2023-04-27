package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"acsp/internal/constants"
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

func (c ContestsDatabase) GetByID(ctx context.Context, contestID int) (model.Contest, error) {
	var contest model.Contest

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		constants.ContestsTable)

	err := c.db.Get(&contest, query, contestID)
	if err != nil {
		return model.Contest{}, errors.Wrap(err, "error when executing the query")
	}

	return contest, nil
}

func (c ContestsDatabase) GetAll(ctx context.Context) ([]model.Contest, error) {
	var contests []model.Contest

	query := fmt.Sprintf("SELECT * FROM %s",
		constants.ContestsTable)

	err := c.db.Select(&contests, query)
	if err != nil {
		return nil, errors.Wrap(err, "error when getting the contest")
	}

	return contests, nil
}
