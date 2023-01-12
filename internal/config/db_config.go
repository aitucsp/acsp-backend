package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"acsp/internal/constants"
)

func NewClientPostgres(ctx context.Context, cancel context.CancelFunc, config *PostgresConfig) (*sqlx.DB, error) {
	connectionQuery := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database)

	db, err := sqlx.Open(constants.DatabaseName, connectionQuery)
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
