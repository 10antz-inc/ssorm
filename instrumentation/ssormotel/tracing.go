package ssormotel

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/codes"
)

type (
	Read func(ctx context.Context) error
	Write func(ctx context.Context) (int64, error)
)

type Tracing interface {
	StartForRead(context.Context, Read) error
	StartForWrite(context.Context, Write) (int64, error)
	SetStatement(string)
	UnsetStatement()
}

type tracing struct {
	conf *config

	spanOpts []trace.SpanStartOption
}

func NewTracing(opts ...Option) *tracing {
	conf := newConfig(opts...)

	return &tracing{
		conf: conf,

		spanOpts: []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(conf.attrs...),
		},
	}
}

func (t *tracing) SetStatement(statement string) {
	t.conf.statement = statement
}

func (t *tracing) UnsetStatement() {
	t.conf.statement = ""
}

func (t *tracing) StartForRead(ctx context.Context, f Read) error {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return f(ctx)
	}

	spanOpts := t.spanOpts
	if t.isEnableStatement() {
		spanOpts = append(spanOpts, trace.WithAttributes(semconv.DBStatementKey.String(t.conf.statement)))
	}

	spanCtx, span := t.conf.tracer.Start(ctx, t.makeSpanName(), spanOpts...)
	defer func() {
		t.UnsetStatement()
		span.End()
	}()

	if err := f(ctx); err != nil {
		recordError(spanCtx, span, err)
		return err
	}
	return nil
}

func (t *tracing) StartForWrite(ctx context.Context, f Write) (int64, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return f(ctx)
	}

	spanOpts := t.spanOpts
	if t.isEnableStatement() {
		spanOpts = append(spanOpts, trace.WithAttributes(semconv.DBStatementKey.String(t.conf.statement)))
	}

	ctx, span := t.conf.tracer.Start(ctx, t.makeSpanName(), spanOpts...)
	defer func() {
		t.UnsetStatement()
		span.End()
	}()

	row, err := f(ctx)
	if err != nil {
		recordError(ctx, span, err)
		return row, err
	}
	return row, nil
}

func (t *tracing) isEnableStatement() bool {
	return t.conf.enableQueryStatement && t.conf.statement != ""
}

func (t *tracing) makeSpanName() string {
	pc, _, _, _ := runtime.Caller(3)
	return runtime.FuncForPC(pc).Name()
}

func recordError(ctx context.Context, span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
