package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"acsp/internal/apperror"
	"acsp/internal/config"
	"acsp/internal/logs"
	"acsp/internal/model"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user model.User) (int, error) {
	logs.Log().Info("Creating a user in database...")

	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, email, password) values ($1, $2, $3) RETURNING id",
		config.UsersTable)

	row := r.db.QueryRow(query, user.Name, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		logs.Log().Info(err.Error())
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(ctx context.Context, email, password string) (*model.User, error) {
	logs.Log().Info("Getting a user from database...")
	var user model.User

	query := fmt.Sprintf("SELECT * FROM users WHERE email=$1 LIMIT 1")

	err := r.db.Get(&user, query, email)
	if err != nil {
		logs.Log().Info(err.Error())
		return &model.User{}, apperror.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logs.Log().Info("")
		return &model.User{}, apperror.ErrPasswordMismatch
	}

	return &user, nil
}

func (r *AuthPostgres) GetByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User

	query := fmt.Sprintf("SELECT * FROM users WHERE id=$1")
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logs.Log().Info(err.Error())
		return nil, apperror.ErrUserNotFound
	}

	return &user, nil
}

func (r *AuthPostgres) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User

	query := fmt.Sprintf("SELECT * FROM users WHERE email=$1")
	row := r.db.QueryRow(query, email)

	if err := row.Scan(&user); err != nil {
		return nil, apperror.ErrUserNotFound
	}

	return user, nil
}

func (r *AuthPostgres) GetAll(ctx context.Context) (*[]model.User, error) {
	var users []model.User

	query := fmt.Sprintf("SELECT * FROM users")

	err := r.db.Select(&users, query)
	if err != nil {
		return nil, apperror.ErrUserNotFound
	}

	return &users, err
}
