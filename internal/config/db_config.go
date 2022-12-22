package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"acsp/internal/apperror"
	"acsp/internal/logs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI      string `mapstructure:"mongoURI"`
	Username string `mapstructure:"mongoUsername"`
	Password string `mapstructure:"mongoPassword"`
	Name     string `mapstructure:"databaseName"`
}

type PostgresConfig struct {
	Username string `yaml:"username" env:"PSQL_USERNAME" env-required:"true"`
	Password string `yaml:"password" env:"PSQL_PASSWORD" env-required:"true"`
	Host     string `yaml:"host" env:"PSQL_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"PSQL_PORT" env-required:"true"`
	Database string `yaml:"database" env:"PSQL_DATABASE" env-required:"true"`
}

func NewClientPostgres(ctx context.Context, cancel context.CancelFunc, config PostgresConfig) (*sql.DB, error) {
	connectionQuery := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database)

	db, err := sql.Open("postgres", connectionQuery)
	if err != nil {
		return nil, err
	}

	defer cancel()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MongoNewClient(ctx context.Context, cancel context.CancelFunc, mongoCfg *MongoConfig) (*mongo.Client, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	err := viper.UnmarshalKey("mongo", &mongoCfg)
	if err != nil {
		return nil, err
	}

	opts := options.Client().ApplyURI(mongoCfg.URI)

	if mongoCfg.Username != "" && mongoCfg.Password != "" {
		opts.SetAuth(
			options.Credential{
				Username: mongoCfg.Username,
				Password: mongoCfg.Password,
			})
	} else {
		return nil, apperror.ErrBadCredentials
	}

	logs.Log().Info("Enabling new mongodb client")
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	defer cancel()

	logs.Log().Info("Connecting to the database")
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	logs.Log().Info("Pinging the database")
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
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
