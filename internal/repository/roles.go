package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"acsp/internal/apperror"
	"acsp/internal/constants"
	"acsp/internal/logging"
	"acsp/internal/model"
)

type RolesDatabase struct {
	db *sqlx.DB
}

func NewRolesRepository(db *sqlx.DB) *RolesDatabase {
	return &RolesDatabase{
		db: db,
	}
}

func (r RolesDatabase) CreateRole(ctx context.Context, name string) error {
	l := logging.LoggerFromContext(ctx).With(zap.String("roleName", name))

	query := fmt.Sprintf(
		`INSERT INTO %s (name) values ($1) 
               RETURNING id`,
		constants.RolesTable)

	_, err := r.db.Exec(query, name)
	if err != nil {
		l.Error("Error when creating new role in database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) UpdateRole(ctx context.Context, roleID int, newName string) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("roleID", roleID), zap.String("roleName", newName))

	query := fmt.Sprintf(
		`UPDATE %s (name) 
				SET name = $1 
				WHERE id = $2`,
		constants.RolesTable)

	_, err := r.db.Exec(query, constants.RolesTable, roleID, newName)
	if err != nil {
		l.Error("Error when updating role in database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) DeleteRole(ctx context.Context, roleID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("roleID", roleID))

	query := fmt.Sprintf(
		`DELETE FROM %s 
				WHERE id = $1`,
		constants.RolesTable)

	_, err := r.db.Exec(query, roleID)
	if err != nil {
		l.Error("Error when deleting role from database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) SaveUserRole(ctx context.Context, userID, roleID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", userID), zap.Int("roleID", roleID))

	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, role_id) 
				VALUES ($1, $2) RETURNING id`,
		constants.UserRolesTable)

	_, err := r.db.Exec(query, userID, roleID)
	if err != nil {
		l.Error("Error when saving user role in database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) DeleteUserRole(ctx context.Context, userID, roleID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", userID), zap.Int("roleID", roleID))

	query := fmt.Sprintf(
		`DELETE FROM %s 
				WHERE user_id = $1 AND role_id = $2`,
		constants.UserRolesTable)

	_, err := r.db.Exec(query, constants.UserRolesTable, userID, roleID)
	if err != nil {
		l.Error("Error when deleting user role from database", zap.Error(err))

		return err
	}

	return nil
}

func (r RolesDatabase) GetUserRoles(ctx context.Context, userID int) ([]model.Role, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", userID))

	var roles []model.Role

	query := fmt.Sprintf(
		`SELECT r.* FROM %s u INNER JOIN %s r ON u.role_id = r.id
				WHERE user_id = $1`,
		constants.UserRolesTable, constants.RolesTable)

	err := r.db.Select(&roles, query, userID)
	if err != nil {
		l.Error("Error when getting user roles from database", zap.Error(err))

		return nil, errors.Wrap(apperror.ErrUserNotFound, "User not found")
	}

	return roles, err
}
