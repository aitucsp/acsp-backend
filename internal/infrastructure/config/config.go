package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

// Environment controls environment-dependent features.
type Environment string

const (
	// EnvironmentDevelopment enables development mode.
	EnvironmentDevelopment Environment = "development"
	// EnvironmentProduction enables production mode.
	EnvironmentProduction Environment = "production"
)

func parseEnvironment(s string) (Environment, error) {
	switch Environment(s) {
	case EnvironmentDevelopment, EnvironmentProduction:
		return Environment(s), nil

	default:
		return "", fmt.Errorf("unknown environment: %v", s)
	}
}

// HostConfig controls application-wide behavior.
type HostConfig struct {
	Environment     Environment `envconfig:"ENVIRONMENT"`
	ShutdownTimeout int         `envconfig:"SERVER_SHUTDOWNTIMEOUT"`
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

// ServerConfig type defines server properties in a config
type ServerConfig struct {
	Port        int    `envconfig:"SERVER_PORT"`
	TLSPort     int    `envconfig:"SERVER_TLSPORT"`
	Host        string `envconfig:"SERVER_HOST"`
	Certificate string `envconfig:"SERVER_CERTIFICATE"`
	Key         string `envconfig:"SERVER_KEY"`
}

// DatabaseConfig type defines server properties in a config
type DatabaseConfig struct {
	Port           string `envconfig:"DB_PORT"`
	Host           string `envconfig:"DB_HOST"`
	Name           string `envconfig:"DB_NAME"`
	Username       string `envconfig:"DB_USERNAME"`
	UserPwd        string `envconfig:"DB_USERNAME_PWD"`
	Timeout        int    `envconfig:"TIMEOUT"`
	UserIDHashCost int    `envconfig:"USERID_HASH_COST"`
}

// BaseConfig controls common features.
type BaseConfig struct {
	Host   *HostConfig
	DB     *DatabaseConfig
	Logger *LoggerConfig
}

// LoggerConfig controls logger behavior.
type LoggerConfig struct {
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
		return nil, err
	}

	ess := []string(nil)
	err = json.Unmarshal([]byte(p.Get(prefix+"_ERRORSINKS", `["stderr"]`)), &ess)
	if err != nil {
		return nil, err
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

// RequestLoggingConfig controls request logging.
type RequestLoggingConfig struct {
	Enable   bool `envconfig:"REQUESTLOGGING_ENABLE"`
	DumpBody bool `envconfig:"REQUESTLOGGING_DUMPBODY"`
}

func newRequestLoggingConfig(p Provider) *RequestLoggingConfig {
	const prefix = "REQUESTLOGGING"

	return &RequestLoggingConfig{
		Enable:   getBool(p, prefix+"_ENABLE", false),
		DumpBody: getBool(p, prefix+"_DUMPBODY", false),
	}
}

// Provider represents a configuration store backed by a key-value mapping.
type Provider interface {
	Get(key, fallback string) string
}

type dotenvProvider struct {
	values map[string]string
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

func (p *dotenvProvider) Get(key, fallback string) string {
	v, ok := p.values[key]
	if ok {
		return v
	}

	return fallback
}

type providerChain struct {
	providers []Provider
}

// NewProviderChain allows value overriding by chaining multiple Provider.
func NewProviderChain(ps ...Provider) Provider {
	return &providerChain{
		providers: ps,
	}
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

// NewBaseConfig creates a BaseConfig.
func NewBaseConfig(p Provider) (*BaseConfig, error) {
	c := BaseConfig{
		DB: &DatabaseConfig{
			Port:           p.Get("DB_PORT", "3306"),
			Host:           p.Get("DB_HOST", "localhost"),
			Name:           p.Get("DB_NAME", "spbibot"),
			Username:       p.Get("DB_USERNAME", "user"),
			UserPwd:        p.Get("DB_USERNAME_PWD", ""),
			Timeout:        getInt(p, "TIMEOUT", 0),
			UserIDHashCost: getInt(p, "USERID_HASH_COST", 10),
		},
	}

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

	return &c, nil
}

func getInt(p Provider, key string, fallback int) int {
	v := p.Get(key, "")
	i, err := strconv.Atoi(v)
	if err == nil {
		return i
	}

	return fallback
}

func getBool(p Provider, key string, fallback bool) bool {
	v := p.Get(key, strconv.FormatBool(fallback))
	b, err := strconv.ParseBool(v)
	if err == nil {
		return b
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
