package db

import (
	"acsp/internal/config"
)

type PostgresConfig struct {
	Username string `yaml:"username" env:"PSQL_USERNAME"`
	Password string `yaml:"password" env:"PSQL_PASSWORD"`
	Host     string `yaml:"host" env:"PSQL_HOST"`
	Port     string `yaml:"port" env:"PSQL_PORT"`
	Database string `yaml:"database" env:"PSQL_DATABASE"`
}

func LoadPostgresConfig(path config.Provider) (*PostgresConfig, error) {
	config := &PostgresConfig{
		Username: path.Get("PSQL_USERNAME", ""),
		Password: path.Get("PSQL_PASSWORD", ""),
		Host:     path.Get("PSQL_HOST", ""),
		Port:     path.Get("PSQL_PORT", ""),
		Database: path.Get("PSQL_DATABASE", ""),
	}

	return config, nil
}
