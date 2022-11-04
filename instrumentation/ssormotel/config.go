package ssormotel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

const (
	name = "github.com/10antz-inc/ssorm/ssormotel"
)

type config struct {
	tp     trace.TracerProvider
	tracer trace.Tracer

	attrs []attribute.KeyValue
}

type Option interface {
	apply(conf *config)
}

type option func(conf *config)

func (fn option) apply(conf *config) {
	fn(conf)
}

func newConfig(opts ...Option) *config {
	tp := otel.GetTracerProvider()
	conf := &config{
		tp: tp,
		tracer: tp.Tracer(name),
		attrs: []attribute.KeyValue{
			semconv.DBSystemKey.String("spanner"),
		},
	}
	for _, opt := range opts {
		opt.apply(conf)
	}
	return conf
}

func WithAttributes(attrs ...attribute.KeyValue) Option {
	return option(func(conf *config) {
		conf.attrs = append(conf.attrs, attrs...)
	})
}

func WithConnectName(name string) Option {
	return option(func(conf * config) {
		semconv.DBConnectionStringKey.String(name)
	})
}

func WithTracerProvider(provider trace.TracerProvider) Option {
	return option(func(conf *config) {
		conf.tp = provider
	})
}
