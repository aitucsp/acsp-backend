package db

import "github.com/spf13/viper"

type RedisConfig struct {
	Addr     string `yaml:"address" env:"REDIS_ADDRESS"`
	Database int    `yaml:"database" env:"REDIS_DATABASE"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	Port     string `yaml:"port" env:"REDIS_PORT"`
}

func LoadRedisConfig(path string) (*RedisConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile("base.env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &RedisConfig{
		Addr:     viper.GetString("REDIS_ADDRESS"),
		Database: viper.GetInt("REDIS_DATABASE"),
		Password: viper.GetString("REDIS_PASSWORD"),
		Port:     viper.GetString("REDIS_PORT"),
	}

	return config, err
}
