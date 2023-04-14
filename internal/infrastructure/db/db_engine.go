package db

import (
	"context"
	"crypto/tls"
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

func NewDBEngine(db *sqlx.DB, c redis.Client) *DBEngine {
	return &DBEngine{
		DB:    db,
		Cache: c,
	}
}

// NewDBClient creates a new database connection
func NewDBClient(ctx context.Context, cancel context.CancelFunc, c *PostgresConfig) (*sqlx.DB, error) {
	connectionQuery := fmt.Sprintf(
		`host=%s 
				port=%s 
				user=%s 
				password=%s 
				dbname=%s 
				sslmode=%s`,
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSL)

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

// NewClientRedis creates a new redis connection
func NewClientRedis(ctx context.Context, cancel context.CancelFunc, config *RedisConfig) (*redis.Client, error) {
	URL, err := redis.ParseURL(config.Addr)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(
		&redis.Options{
			Addr:         URL.Addr,
			DB:           URL.DB,
			Password:     config.Password,
			DialTimeout:  time.Second * 10,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		},
	)

	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
