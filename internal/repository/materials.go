package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"acsp/internal/constants"
	"acsp/internal/logging"
	"acsp/internal/model"
)

type MaterialsDatabase struct {
	db *sqlx.DB
}

func NewMaterialsRepository(db *sqlx.DB) *MaterialsDatabase {
	return &MaterialsDatabase{
		db: db,
	}
}

func (m *MaterialsDatabase) Create(ctx context.Context, material model.Material) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("materialID", material.ID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("INSERT INTO %s (user_id, topic, description) VALUES ($1, $2, $3) RETURNING id",
		constants.MaterialsTable)

	stmt, err := m.db.PrepareContext(ctx, query)
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

	res, err := stmt.ExecContext(ctx, material.Author.ID, material.Topic, material.Description)
	if err != nil {
		l.Error("Error when creating the material in database", zap.Error(err))

		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when creating the material in database", zap.Error(err))

		return errors.Wrap(err, "error when creating the material in database")
	}

	return nil
}

func (m *MaterialsDatabase) Update(ctx context.Context, material model.Material) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("materialID", material.ID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("UPDATE %s SET topic = $1, description = $2, updated_at = now() WHERE id = $3 AND user_id = $4",
		constants.MaterialsTable)

	stmt, err := m.db.PrepareContext(ctx, query)
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

	res, err := stmt.ExecContext(ctx, query, material.Topic, material.Description, material.ID, material.Author.ID)
	if err != nil {
		l.Error("Error when update the article in database", zap.Error(err))

		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when creating the material in database", zap.Error(err))

		return errors.Wrap(err, "error when creating the material in database")
	}

	return nil
}

func (m *MaterialsDatabase) Delete(ctx context.Context, userID int, materialID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("materialID", materialID), zap.Int("userID", userID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", constants.MaterialsTable)

	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "error when preparing the query")
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			l.Error("Error when closing the statement", zap.Error(err))
		}
	}(stmt)

	res, err := stmt.ExecContext(ctx, query, materialID, userID)
	if err != nil {
		return errors.Wrap(err, "error when executing the query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		l.Error("Error when getting the rows affected", zap.Error(err))

		return errors.Wrap(err, "error when getting the rows affected")
	}

	if rowsAffected < 1 {
		l.Error("Error when creating the material in database", zap.Error(err))

		return errors.Wrap(err, "error when creating the material in database")
	}

	return nil
}

func (m *MaterialsDatabase) GetByID(ctx context.Context, materialID int) (model.Material, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("materialID", materialID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var material model.Material

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		constants.MaterialsTable)

	err := m.db.GetContext(ctx, &material, query, materialID)
	if err != nil {
		l.Error("Error when getting the material", zap.Error(err))

		return model.Material{}, errors.Wrap(err, "error when executing the query")
	}

	return material, nil
}

func (m *MaterialsDatabase) GetAllByUserID(ctx context.Context, userID int) ([]model.Material, error) {
	var materials []model.Material

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		constants.MaterialsTable)

	err := m.db.Select(&materials, query, userID)
	if err != nil {
		return []model.Material{}, errors.Wrap(err, "error when getting the materials")
	}

	return materials, nil
}

func (m *MaterialsDatabase) GetAll(ctx context.Context) ([]model.Material, error) {
	var materials []model.Material

	query := fmt.Sprintf("SELECT * FROM %s",
		constants.MaterialsTable)

	err := m.db.Select(&materials, query)
	if err != nil {
		return nil, errors.Wrap(err, "error when getting the materials")
	}

	return materials, nil
}
