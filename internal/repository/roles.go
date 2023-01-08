package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"acsp/internal/apperror"
	"acsp/internal/config"
	"acsp/internal/logging"
	"acsp/internal/model"
)

type RolesDatabase struct {
	db *sqlx.DB
}

func NewRolesPostgres(db *sqlx.DB) *RolesDatabase {
	return &RolesDatabase{
		db: db,
	}
}

func (r RolesDatabase) CreateRole(ctx context.Context, name string) error {
	l := logging.LoggerFromContext(ctx).With(zap.String("roleName", name))

	query := fmt.Sprintf("INSERT INTO %s (name) values ($1) RETURNING id", config.RolesTable)

	_, err := r.db.Exec(query, name)
	if err != nil {
		l.Error("Error when creating new role in database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) UpdateRole(ctx context.Context, roleID int, newName string) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("roleID", roleID), zap.String("roleName", newName))

	query := fmt.Sprintf("UPDATE %s (name) SET name = $1 WHERE id = $2",
		config.RolesTable)

	_, err := r.db.Exec(query, config.RolesTable, roleID, newName)
	if err != nil {
		l.Error("Error when updating role in database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) DeleteRole(ctx context.Context, roleID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("roleID", roleID))

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", config.RolesTable)

	_, err := r.db.Exec(query, roleID)
	if err != nil {
		l.Error("Error when deleting role from database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) SaveUserRole(ctx context.Context, userID, roleID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", userID), zap.Int("roleID", roleID))

	query := fmt.Sprintf("INSERT INTO %s (user_id, role_id) values ($1, $2) RETURNING id",
		config.UserRolesTable)

	_, err := r.db.Exec(query, userID, roleID)
	if err != nil {
		l.Error("Error when saving user role in database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) DeleteUserRole(ctx context.Context, userID, roleID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", userID), zap.Int("roleID", roleID))

	query := fmt.Sprint("DELETE FROM $1 WHERE user_id = $2 AND role_id = $3")

	_, err := r.db.Exec(query, config.UserRolesTable, userID, roleID)
	if err != nil {
		l.Error("Error when deleting user role from database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) GetUserRoles(ctx context.Context, userID int) ([]model.Role, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", userID))

	var roles []model.Role

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", config.UserRolesTable)

	err := r.db.Select(&roles, query, userID)
	if err != nil {
		l.Error("Error when getting user roles from database", zap.Error(err))

		return nil, errors.Wrap(apperror.ErrUserNotFound, "user not found")
	}

	return roles, err
}
