package db

import (
	"acsp/internal/config"
)

// RedisConfig is the configuration for the redis database
type RedisConfig struct {
	Addr     string `yaml:"address" env:"REDIS_ADDRESS"`
	Database int    `yaml:"database" env:"REDIS_DATABASE"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
}

// LoadRedisConfig loads the redis configuration
func LoadRedisConfig(p config.Provider) (*RedisConfig, error) {
	c := &RedisConfig{
		Addr:     p.Get("REDIS_ADDRESS", ""),
		Password: p.Get("REDIS_PASSWORD", ""),
		Database: config.GetInt(p, "REDIS_DATABASE", 0),
	}

	return c, nil
}
