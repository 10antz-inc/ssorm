package ssorm

import (
	"io"
	"time"
	"context"

	"github.com/sirupsen/logrus"
)

type ILogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	WithContext(ctx context.Context) *logrus.Entry
}

func NewLogger(out io.Writer) ILogger {
	return &logrus.Logger{
		Out: out,
		Formatter: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "severity",
				logrus.FieldKeyMsg:   "message",
			},
			TimestampFormat: time.RFC3339Nano,
		},
		Level: logrus.DebugLevel,
	}
}
