package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"

	"acsp/internal/constants"
)

type DBEngine struct {
	DB    *sqlx.DB
	Cache redis.Client
}

func NewDBEngine(db *sqlx.DB, cache redis.Client) *DBEngine {
	return &DBEngine{
		DB:    db,
		Cache: cache,
	}
}
func NewDBClient(ctx context.Context, cancel context.CancelFunc, config *PostgresConfig) (*sqlx.DB, error) {
	connectionQuery := fmt.Sprintf(
		`host=%s 
				port=%s 
				user=%s 
				password=%s 
				dbname=%s 
				sslmode=disable`,
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

func NewClientRedis(ctx context.Context, cancel context.CancelFunc, config *RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(
		&redis.Options{
			Addr:     config.Addr,
			DB:       config.Database,
			Password: config.Password,
		},
	)

	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
