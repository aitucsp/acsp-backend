package logging

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"acsp/internal/config"
)

type timeLayout string

const (
	layoutISO8601 timeLayout = "2006-01-02T15:04:05.000Z0700"
)

type key string

const (
	keyTime key = "time"
)

// Builder configures a zap.Logger.
type Builder struct {
	hostConfig     *config.HostConfig
	loggerConfig   *config.LoggerConfig
	fallbackLogger *log.Logger
}

// NewBuilder creates a Builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// WithHostConfig adds a config.HostConfig.
func (b *Builder) WithHostConfig(h *config.HostConfig) *Builder {
	b.hostConfig = h

	return b
}

// WithLoggerConfig adds a config.LoggerConfig.
func (b *Builder) WithLoggerConfig(l *config.LoggerConfig) *Builder {
	b.loggerConfig = l

	return b
}

// WithFallbackLogger adds a log.Logger.
func (b *Builder) WithFallbackLogger(l *log.Logger) *Builder {
	b.fallbackLogger = l

	return b
}

// NewLogger creates a zap.Logger.
func (b *Builder) NewLogger() (*zap.Logger, error) {
	loggerConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(b.loggerConfig.Level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          b.loggerConfig.Encoding,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:       "message",
			LevelKey:         "level",
			TimeKey:          string(keyTime),
			NameKey:          "name",
			CallerKey:        "caller",
			FunctionKey:      zapcore.OmitKey,
			StacktraceKey:    "stacktrace",
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      b.loggerConfig.LevelEncoder,
			EncodeTime:       zapcore.TimeEncoderOfLayout(string(layoutISO8601)),
			EncodeDuration:   zapcore.StringDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			EncodeName:       zapcore.FullNameEncoder,
			ConsoleSeparator: "\t",
		},
		OutputPaths:      b.loggerConfig.Sinks,
		ErrorOutputPaths: b.loggerConfig.ErrorSinks,
		InitialFields:    nil,
	}
	switch b.hostConfig.Environment {
	case config.EnvironmentDevelopment:
		loggerConfig.Development = true

	case config.EnvironmentProduction:
		loggerConfig.Development = false

	default:
		return nil, fmt.Errorf("unknown environment: %v", b.hostConfig.Environment)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	o := pathOptions{
		Host: hostname,
	}

	err = zap.RegisterSink("lumberjack", lumberjackSinkFactory(b.loggerConfig, &o))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't register sink")
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't build logger")
	}

	zap.ReplaceGlobals(logger)

	_, err = zap.RedirectStdLogAt(logger.Named("std"), zapcore.InfoLevel)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't redirect std logger")
	}

	return logger, nil
}

type ctxLogger struct{}

// ContextWithLogger adds logger to context
func ContextWithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

// LoggerFromContext returns logger from context
func LoggerFromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*zap.Logger); ok {
		return l
	}
	return zap.L()
}
