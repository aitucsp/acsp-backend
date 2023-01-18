package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"

	"acsp/internal/apperror"
)

const (
	EnvLocal = "local"
)

const (
	// EnvironmentDevelopment enables development mode.
	EnvironmentDevelopment Environment = "development"
	// EnvironmentProduction enables production mode.
	EnvironmentProduction Environment = "production"
)

type (
	Config struct {
		Environment string
		HTTP        *HTTPConfig
		Auth        *AuthConfig
		Logger      *LoggerConfig
		Host        *HostConfig
	}

	AuthConfig struct {
		JWT JWTConfig
	}

	JWTConfig struct {
		AccessTokenTTL     time.Duration `envconfig:"ACCESS_TOKEN_TTL"`
		RefreshTokenTTL    time.Duration `envconfig:"REFRESH_TOKEN_TTL"`
		AccessTokenSecret  string        `envconfig:"ACCESS_TOKEN_SECRET_KEY"`
		RefreshTokenSecret string        `envconfig:"REFRESH_TOKEN_SECRET_KEY"`
	}

	HTTPConfig struct {
		Host               string        `envconfig:"HOST"`
		Port               string        `envconfig:"PORT"`
		ReadTimeout        time.Duration `envconfig:"READ_TIMEOUT"`
		WriteTimeout       time.Duration `envconfig:"WRITE_TIMEOUT"`
		MaxHeaderMegabytes int           `envconfig:"MAX_HEADER_BYTES"`
	}

	LoggerConfig struct {
		Level        zapcore.Level        `envconfig:"LOGGER_LEVEL"`
		Encoding     string               `envconfig:"LOGGER_ENCODING"`
		LevelEncoder zapcore.LevelEncoder `envconfig:"LOGGER_LEVELENCODER"`
		Sinks        []string             `envconfig:"LOGGER_SINKS"`
		ErrorSinks   []string             `envconfig:"LOGGER_ERRORSINKS"`
		MaxSizeMB    int                  `envconfig:"LOGGER_MAXSIZEMB"`
		MaxAgeDays   int                  `envconfig:"LOGGER_MAXAGEDAYS"`
		MaxBackups   int                  `envconfig:"LOGGER_MAXBACKUPS"`
		BatchSize    uint                 `envconfig:"LOGGER_BATCHSIZE"`
	}

	// Environment controls environment-dependent features.
	Environment string

	// HostConfig controls application-wide behavior.
	HostConfig struct {
		Environment     Environment `envconfig:"ENVIRONMENT"`
		ShutdownTimeout int         `envconfig:"SERVER_SHUTDOWNTIMEOUT"`
	}
)

func Init(configsDir string) (*Config, error) {
	viper.SetConfigFile("base.env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := parseConfigFile(configsDir, viper.GetString("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// NewBaseConfig creates a BaseConfig.
func NewBaseConfig(p Provider) (*Config, error) {
	c := Config{
		Environment: p.Get("ENVIRONMENT", "development"),
		HTTP:        &HTTPConfig{},
	}

	http, err := newHTTPConfig(p)
	if err != nil {
		return nil, err
	}

	c.HTTP = http

	h, err := newHostConfig(p)
	if err != nil {
		return nil, err
	}

	c.Host = h

	l, err := newLoggerConfig(p)
	if err != nil {
		return nil, err
	}

	c.Logger = l

	a, err := newAuthConfig(p)
	if err != nil {
		return nil, err
	}

	c.Auth = a

	return &c, nil
}

func newHostConfig(p Provider) (*HostConfig, error) {
	e, err := parseEnvironment(p.Get("ENV", string(EnvironmentDevelopment)))
	if err != nil {
		return nil, err
	}

	return &HostConfig{
		Environment:     e,
		ShutdownTimeout: getInt(p, "SERVER_SHUTDOWNTIMEOUT", 10),
	}, nil
}

type dotenvProvider struct {
	values map[string]string
}

func (p *dotenvProvider) Get(key, fallback string) string {
	v, ok := p.values[key]
	if ok {
		return v
	}

	return fallback
}

// NewDotenvProvider creates a .env file-backed Provider.
func NewDotenvProvider(filepath string) (Provider, error) {
	vs, err := godotenv.Read(filepath)
	if err != nil {
		return nil, err
	}

	return &dotenvProvider{
		values: vs,
	}, nil
}

type providerChain struct {
	providers []Provider
}

func (c *providerChain) Get(key, fallback string) string {
	for _, p := range c.providers {
		v := p.Get(key, fallback)
		if v != fallback {
			return v
		}
	}

	return fallback
}

// NewProviderChain allows value overriding by chaining multiple Provider.
func NewProviderChain(ps ...Provider) Provider {
	return &providerChain{
		providers: ps,
	}
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	return viper.UnmarshalKey("auth", &cfg.Auth.JWT)
}

func parseConfigFile(folder, fileName string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if fileName == EnvLocal {
		return nil
	}

	viper.SetConfigName(fileName)

	return viper.MergeInConfig()
}

func parseEnvironment(s string) (Environment, error) {
	switch Environment(s) {
	case EnvironmentDevelopment, EnvironmentProduction:
		return Environment(s), nil

	default:
		return "", fmt.Errorf("unknown environment: %v", s)
	}
}

// Provider represents a configuration store backed by a key-value mapping.
type Provider interface {
	Get(key, fallback string) string
}

func newHTTPConfig(p Provider) (*HTTPConfig, error) {
	const prefix = "HTTP"

	return &HTTPConfig{
		Host:               p.Get(prefix+"_HOST", "localhost"),
		Port:               p.Get(prefix+"_PORT", "8080"),
		ReadTimeout:        getDuration(p, prefix+"_READ_TIMEOUT", 10),
		WriteTimeout:       getDuration(p, prefix+"_WRITE_TIMEOUT", 10),
		MaxHeaderMegabytes: getInt(p, prefix+"_MAX_HEADER_BYTES", 1),
	}, nil
}

func newLoggerConfig(p Provider) (*LoggerConfig, error) {
	const prefix = "LOGGER"

	var l zapcore.Level
	err := l.UnmarshalText([]byte(p.Get(prefix+"_LEVEL", "info")))
	if err != nil {
		return nil, err
	}

	var le zapcore.LevelEncoder
	err = le.UnmarshalText([]byte(p.Get(prefix+"_LEVELENCODER", "capitalColor")))
	if err != nil {
		return nil, err
	}

	f := getBool(p, prefix+"_ENABLEFILE", true)
	s := getBool(p, prefix+"_ENABLESTDOUT", true)
	if !f && !s {
		return nil, fmt.Errorf("at least one sink must be enabled")
	}

	ss := []string(nil)
	err = json.Unmarshal([]byte(p.Get(prefix+"_SINKS", `["stdout"]`)), &ss)
	if err != nil {
		return nil, errors.Wrap(err, "no sinks set")
	}

	ess := []string(nil)
	err = json.Unmarshal([]byte(p.Get(prefix+"_ERRORSINKS", `["stderr"]`)), &ess)
	if err != nil {
		return nil, errors.Wrap(err, "no error sinks set")
	}

	return &LoggerConfig{
		Level:        l,
		Encoding:     p.Get("LOGGER_ENCODING", "console"),
		LevelEncoder: le,
		Sinks:        ss,
		ErrorSinks:   ess,
		MaxSizeMB:    getInt(p, "LOGGER_MAXSIZEMB", 128),
		MaxAgeDays:   getInt(p, "LOGGER_MAXAGEDAYS", 168),
		MaxBackups:   getInt(p, "LOGGER_MAXBACKUPS", 16),
		BatchSize:    getUint(p, "LOGGER_BATCHSIZE", 2),
	}, nil
}

func newAuthConfig(p Provider) (*AuthConfig, error) {

	const prefix = "JWT"

	var accessTokenKey string
	accessTokenKey = p.Get(prefix+"_ACCESS_TOKEN_SECRET_KEY", "")
	if accessTokenKey == "" {
		return nil, apperror.ErrEnvVariableParsing
	}

	var refreshTokenKey string
	refreshTokenKey = p.Get(prefix+"_ACCESS_TOKEN_SECRET_KEY", "")
	if refreshTokenKey == "" {
		return nil, apperror.ErrEnvVariableParsing
	}

	var accessTokenTTL time.Duration
	accessTokenTTL = getDuration(p, prefix+"_ACCESS_TOKEN_TTL", 0)
	if accessTokenTTL == 0 {
		return nil, apperror.ErrEnvVariableParsing
	}

	var refreshTokenTTL time.Duration
	refreshTokenTTL = getDuration(p, prefix+"_ACCESS_TOKEN_TTL", 0)
	if refreshTokenTTL == 0 {
		return nil, apperror.ErrEnvVariableParsing
	}

	return &AuthConfig{
		JWT: JWTConfig{
			AccessTokenTTL:     accessTokenTTL,
			RefreshTokenTTL:    refreshTokenTTL,
			AccessTokenSecret:  accessTokenKey,
			RefreshTokenSecret: refreshTokenKey,
		},
	}, nil
}

func getBool(p Provider, key string, fallback bool) bool {
	v := p.Get(key, strconv.FormatBool(fallback))
	b, err := strconv.ParseBool(v)
	if err == nil {
		return b
	}

	return fallback
}

func getInt(p Provider, key string, fallback int) int {
	v := p.Get(key, "")
	i, err := strconv.Atoi(v)
	if err == nil {
		return i
	}

	return fallback
}

func getFloat64(p Provider, key string, fallback float64) float64 {
	v := p.Get(key, "")
	f64, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return f64
	}

	return fallback
}

func getInt64(p Provider, key string, fallback int64) int64 {
	v := p.Get(key, "")
	i64, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		return i64
	}

	return fallback
}

func getUint(p Provider, key string, fallback uint) uint {
	v := p.Get(key, "")
	u, err := strconv.ParseUint(v, 10, 0)
	if err == nil {
		return uint(u)
	}

	return fallback
}

func getDuration(p Provider, key string, fallback time.Duration) time.Duration {
	v := p.Get(key, "")
	d, err := time.ParseDuration(v)
	if err == nil {
		return d
	}

	return fallback
}
