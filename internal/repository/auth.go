package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

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

func (r *AuthPostgres) CreateUser(ctx context.Context, user model.User) error {
	l := logging.LoggerFromContext(ctx).With(zap.String("email", user.Email))

	// Begin a transaction to create a user and a role for the user
	tx, err := r.db.Begin()
	if err != nil {
		l.Error("Error when beginning the transaction", zap.Error(err))

		return errors.Wrap(err, "Error when beginning the transaction")
	}

	var userID int
	query := fmt.Sprintf(
		`INSERT INTO %s (name, email, password) 
			    VALUES ($1, $2, $3) 
  				RETURNING id`,
		constants.UsersTable)

	// Create a user
	row := tx.QueryRow(query, user.Name, user.Email, user.Password)

	// Get the user id
	err = row.Scan(&userID)
	if err != nil {
		l.Error("Error when creating user in database", zap.Error(err))

		return errors.Wrap(err, "Error when creating user in database")
	}

	// Rollback the transaction if the user id is less than 1
	if userID < 1 {
		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}
	}

	querySecond := fmt.Sprintf(`INSERT INTO %s (user_id, role_id) 
										VALUES ($1, $2) RETURNING id`, constants.UserRolesTable)
	var userRoleID int

	// Create a role for the user
	err = tx.QueryRow(querySecond, userID, constants.DefaultUserRoleID).Scan(&userRoleID)
	if err != nil {
		l.Error("Error when executing the query", zap.Error(err))

		// Rollback the transaction if the user role id is less than 1
		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back the transaction")
		}

		// Return the error
		return errors.Wrap(err, "Error when executing the query")
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		// Rollback the transaction if the commit fails
		err := tx.Rollback()
		if err != nil {
			return errors.Wrap(apperror.ErrRollback, "Error when rolling back")
		}

		// Return the error
		return errors.Wrap(err, "Error when committing the transaction")
	}

	// Return nil if the transaction is successful
	return nil
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
	query := fmt.Sprintf(`SELECT 
									u.id,
									u.email,
									u.name,
									u.password,
									u.created_at,
									u.updated_at,
									u.is_admin,
									ARRAY_AGG(r.name) AS roles,
									u.image_url
								FROM %s u WHERE id=$1`,
		constants.UsersTable)

	err := r.db.
		QueryRow(query, id).
		Scan(&user.ID,
			&user.Email,
			&user.Name,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.IsAdmin,
			&user.Roles,
			&user.ImageURL,
		)
	if err != nil {
		l.Error("Error getting user by id from database", zap.Error(err))

		return nil, apperror.ErrUserNotFound
	}

	return &user, nil
}

func (r *AuthPostgres) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	query := fmt.Sprintf(`SELECT 
									u.id,
									u.email,
									u.name,	
									u.password,	
									u.created_at,	
									u.updated_at,
									u.is_admin,
									ARRAY_AGG(r.name) AS roles,
									u.image_url
								FROM %s u WHERE email=$1`,
		constants.UsersTable)
	row := r.db.QueryRow(query, email)

	err := row.Scan(&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsAdmin,
		&user.Roles,
		&user.ImageURL)
	if err != nil {
		return nil, errors.Wrap(apperror.ErrEmailNotFound, "email not found")
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

func (r *AuthPostgres) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	l := logging.LoggerFromContext(ctx)

	var isExists bool

	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE email=$1)`, constants.UsersTable)

	row := r.db.QueryRow(query, email)

	err := row.Scan(&isExists)
	if err != nil {
		l.Error("Error when finding user in database", zap.Error(err))

		return false, err
	}

	return isExists, nil
}

func (r *AuthPostgres) UpdateImageURL(ctx context.Context, userID int) error {
	l := logging.LoggerFromContext(ctx).With(zap.Int("userID", userID))

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`UPDATE %s SET image_url = $1, updated_at = now() WHERE id = $2`,
		constants.UsersTable)

	stmt, err := r.db.PrepareContext(ctx, query)
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

	res, err := stmt.Exec("/"+strconv.Itoa(userID), userID)
	if err != nil {
		l.Error("Error when update the user's image url in database", zap.Error(err))

		return errors.Wrap(err, "error when executing query")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error when get rows affected")
	}

	if rows == 0 {
		return errors.Wrap(apperror.ErrWhenUpdatingImageURL, "error when updating a user image url")
	}

	return nil
}
