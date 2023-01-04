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

	query := fmt.Sprint("INSERT INTO $1 (name) values ($2) RETURNING id")

	_, err := r.db.Exec(query, config.RolesTable, name)
	if err != nil {
		l.Error("Error when creating new role in database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) UpdateRole(ctx context.Context, roleID int, newName string) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("roleID", roleID), zap.String("roleName", newName))

	query := fmt.Sprint("UPDATE $1 (name) SET name = $2 WHERE id = $3")

	_, err := r.db.Exec(query, config.RolesTable, roleID, newName)
	if err != nil {
		l.Error("Error when updating role in database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) DeleteRole(ctx context.Context, roleID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("roleID", roleID))

	query := fmt.Sprint("DELETE FROM $1 WHERE id = $2")

	_, err := r.db.Exec(query, config.RolesTable, roleID)
	if err != nil {
		l.Error("Error when deleting role from database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) SaveUserRole(ctx context.Context, userID, roleID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", userID), zap.Int("roleID", roleID))

	query := fmt.Sprint("INSERT INTO $1 (user_id, role_id) values ($2, $3) RETURNING id")

	_, err := r.db.Exec(query, config.UserRolesTable, userID, roleID)
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
		l.Error("Error when saving user role in database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) GetUserRoles(ctx context.Context, userID int) ([]model.Role, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", userID))

	var roles []model.Role

	query := fmt.Sprint("SELECT * FROM $1 WHERE user_id = $2")

	err := r.db.Select(&roles, query, config.UserRolesTable, userID)
	if err != nil {
		l.Error("Error when getting user roles from database", zap.Error(err))

		return nil, errors.Wrap(apperror.ErrUserNotFound, "user not found")
	}

	return roles, err
}
