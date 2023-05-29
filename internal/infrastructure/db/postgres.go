package db

import (
	"acsp/internal/config"
)

// PostgresConfig is the configuration for the postgres database
type PostgresConfig struct {
	Username string `yaml:"username" env:"PSQL_USERNAME"`
	Password string `yaml:"password" env:"PSQL_PASSWORD"`
	Host     string `yaml:"host" env:"PSQL_HOST"`
	Port     string `yaml:"port" env:"PSQL_PORT"`
	Database string `yaml:"database" env:"PSQL_DATABASE"`
	SSL      string `yaml:"ssl" env:"PSQL_SSLMODE"`
}

// LoadPostgresConfig loads the postgres configuration
func LoadPostgresConfig(p config.Provider) (*PostgresConfig, error) {
	c := &PostgresConfig{
		Username: p.Get("PSQL_USERNAME", ""),
		Password: p.Get("PSQL_PASSWORD", ""),
		Host:     p.Get("PSQL_HOST", ""),
		Port:     p.Get("PSQL_PORT", ""),
		Database: p.Get("PSQL_DATABASE", ""),
		SSL:      p.Get("PSQL_SSLMODE", ""),
	}

	return c, nil
}
