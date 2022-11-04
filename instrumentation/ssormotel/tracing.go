package ssormotel

import (
	"go.opentelemetry.io/otel/trace"
)

type tracing struct {
	conf *config

	spanOpts []trace.SpanStartOption
}

func newTracing(opts []Option) *tracing {
	conf := newConfig(opts...)

	return &tracing{
		conf: conf,

		spanOpts: []trace.SpanStartOption{
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(conf.attrs...),
		},
	}
}
