package config

import "github.com/spf13/viper"

type PostgresConfig struct {
	Username string `yaml:"username" env:"PSQL_USERNAME"`
	Password string `yaml:"password" env:"PSQL_PASSWORD"`
	Host     string `yaml:"host" env:"PSQL_HOST"`
	Port     string `yaml:"port" env:"PSQL_PORT"`
	Database string `yaml:"database" env:"PSQL_DATABASE"`
}

func LoadConfig(path string) (*PostgresConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile("base.env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &PostgresConfig{
		Username: viper.GetString("PSQL_USERNAME"),
		Password: viper.GetString("PSQL_PASSWORD"),
		Host:     viper.GetString("PSQL_HOST"),
		Port:     viper.GetString("PSQL_PORT"),
		Database: viper.GetString("PSQL_DATABASE"),
	}

	return config, err
}
