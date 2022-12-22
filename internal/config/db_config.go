package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	UsersTable    = "users"
	ArticlesTable = "articles"
)

func NewClientPostgres(ctx context.Context, cancel context.CancelFunc, config *PostgresConfig) (*sqlx.DB, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("base")
	viper.SetConfigType("env")

	connectionQuery := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database)

	db, err := sqlx.Open("postgres", connectionQuery)
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func IsDuplicate(err error) bool {
	var e mongo.WriteException

	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}

	return false
}
