package ssormotel

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

type (
	Read func(ctx context.Context) error
	Write func(ctx context.Context) (int64, error)
)

type Tracing interface {
	StartForRead(context.Context, Read) error
	StartForWrite(context.Context, Write) (int64, error)
	AddStatementToSpanAttribute(string)
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

func (t *tracing) StartForRead(ctx context.Context, f Read) error {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return f(ctx)
	}

	return func() error {
		ctx, span := t.conf.tracer.Start(ctx, t.makeSpanName(), t.spanOpts...)
		defer span.End()
		return f(ctx)
	}()
}

func (t *tracing) StartForWrite(ctx context.Context, f Write) (int64, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return f(ctx)
	}

	return func() (int64, error) {
		ctx, span := t.conf.tracer.Start(ctx, t.makeSpanName(), t.spanOpts...)
		defer span.End()
		return f(ctx)
	}()
}

func (t *tracing) AddStatementToSpanAttribute(statement string) {
	if !t.conf.enableQueryStatement {
		return
	}

	v := semconv.DBStatementKey.String(statement)
	t.spanOpts = append(t.spanOpts, trace.WithAttributes(attribute.KeyValue{Key: v.Key, Value: v.Value}))
}

func (t *tracing) makeSpanName() string {
	pc, _, _, _ := runtime.Caller(3)
	return runtime.FuncForPC(pc).Name()
}
