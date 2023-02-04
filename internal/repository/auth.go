package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"acsp/internal/apperror"
	"acsp/internal/constants"
	"acsp/internal/logging"
	"acsp/internal/model"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user model.User) (int, error) {
	l := logging.LoggerFromContext(ctx).With(zap.String("email", user.Email))

	var userID int
	query := fmt.Sprintf(
		`INSERT INTO %s (name, email, password) 
			    values ($1, $2, $3) 
  				RETURNING id`,
		constants.UsersTable)

	row := r.db.QueryRow(query, user.Name, user.Email, user.Password)

	err := row.Scan(&userID)
	if err != nil {
		l.Error("Error when creating user in database", zap.Error(err))

		return -1, err
	}

	return userID, nil
}

func (r *AuthPostgres) GetUser(ctx context.Context, email, password string) (*model.User, error) {
	l := logging.LoggerFromContext(ctx).With(zap.String("email", email))

	var user model.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE email=$1 LIMIT 1`,
		constants.UsersTable)

	err := r.db.Get(&user, query, email)
	if err != nil {
		l.Error("Error when getting the user from database", zap.Error(err))

		return &model.User{}, apperror.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		l.Error("Error when comparing passwords of user", zap.Error(err))

		return &model.User{}, apperror.ErrPasswordMismatch
	}

	return &user, nil
}

func (r *AuthPostgres) GetByID(ctx context.Context, id int) (*model.User, error) {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", id))
	var user model.User
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id=$1`, constants.UsersTable)

	err := r.db.
		QueryRow(query, id).
		Scan(&user.ID,
			&user.Email,
			&user.Name,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Roles)
	if err != nil {
		l.Error("Error getting user by id from database", zap.Error(err))

		return nil, apperror.ErrUserNotFound
	}

	return &user, nil
}

func (r *AuthPostgres) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	l := logging.LoggerFromContext(ctx)

	var user model.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE email=$1`, constants.UsersTable)
	row := r.db.QueryRow(query, email)

	if err := row.Scan(&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Roles); err != nil {
		l.Error("Error", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r *AuthPostgres) GetAll(ctx context.Context) (*[]model.User, error) {
	l := logging.LoggerFromContext(ctx)

	var users []model.User

	query := fmt.Sprintf(`SELECT * FROM %s`, constants.UsersTable)

	err := r.db.Select(&users, query)
	if err != nil {
		l.Error("Error when getting users from database", zap.Error(err))

		return nil, errors.Wrap(apperror.ErrUserNotFound, "user not found")
	}

	return &users, err
}

func (r *AuthPostgres) ExistsUserByID(ctx context.Context, id int) (bool, error) {
	l := logging.LoggerFromContext(ctx)

	var isExists bool

	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)`, constants.UsersTable)

	row := r.db.QueryRow(query, id)

	err := row.Scan(&isExists)
	if err != nil {
		l.Error("Error when finding user in database", zap.Error(err))

		return false, err
	}

	return isExists, nil
}
