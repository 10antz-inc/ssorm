package logger

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type Logger struct {
	conf *config
}

func NewLogger(opts ...Option) *Logger {
	return &Logger{
		conf: newConfig(opts...),
	}
}

func (l *Logger) ctx(ctx context.Context) *zerolog.Logger {
	var (
		zctx zerolog.Context = log.Ctx(ctx).With()
	)

	return lo.ToPtr(zctx.Logger())
}

func (l *Logger) ReadLog(ctx context.Context, format string, v ...any) {
	if l.conf.outputQueryLogType.AllowReadLog() {
		l.ctx(ctx).Info().Fields(l.conf.fields).Msgf(format, v...)
	}
}

func (l *Logger) WriteLog(ctx context.Context, format string, v ...any) {
	if l.conf.outputQueryLogType.AllowWriteLog() {
		l.ctx(ctx).Info().Fields(l.conf.fields).Msgf(format, v...)
	}
}

func (l *Logger) ErrorLog(ctx context.Context, err error, format string, v ...any) {
	l.ctx(ctx).Error().Fields(l.conf.fields).AnErr("error", err).Msgf(format, v...)
}
