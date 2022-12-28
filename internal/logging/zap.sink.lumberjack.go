package logging

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"text/template"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"

	"acsp/internal/config"
)

type lumberjackSink struct {
	*lumberjack.Logger
}

func newLumberjackSink(l *lumberjack.Logger) *lumberjackSink {
	return &lumberjackSink{
		Logger: l,
	}
}

func (*lumberjackSink) Sync() error {
	return nil
}

func lumberjackSinkFactory(c *config.LoggerConfig, o *pathOptions) func(u *url.URL) (zap.Sink, error) {
	return func(u *url.URL) (zap.Sink, error) {
		u, err := o.useWith(u)
		if err != nil {
			return nil, err
		}

		if u.Host != "localhost" {
			return nil, fmt.Errorf("host must be localhost")
		}

		l := lumberjack.Logger{
			Filename:   strings.TrimPrefix(u.Path, "/"),
			MaxSize:    c.MaxSizeMB,
			MaxAge:     c.MaxAgeDays,
			MaxBackups: c.MaxBackups,
			LocalTime:  false,
			Compress:   false,
		}
		sink := newLumberjackSink(&l)

		return sink, nil
	}
}

type pathOptions struct {
	Host string
}

func (o *pathOptions) useWith(u *url.URL) (*url.URL, error) {
	source := u.String()
	source, err := url.PathUnescape(source)
	if err != nil {
		return nil, err
	}

	sourceTemplate, err := template.New("").Option("missingkey=error").Parse(source)
	if err != nil {
		return nil, err
	}

	targetBuffer := bytes.Buffer{}
	err = sourceTemplate.Execute(&targetBuffer, *o)
	if err != nil {
		return nil, err
	}

	target := targetBuffer.String()

	return url.Parse(target)
}
