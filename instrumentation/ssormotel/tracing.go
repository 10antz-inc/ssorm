package ssormotel

import (
	// "context"
	// "fmt"

	// "github.com/10antz-inc/ssorm"

	// "go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/codes"
	// semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
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
//
// func (th *tracingHook) DialHook(hook redis.DialHook) redis.DialHook {
// 	return func(ctx context.Context, network, addr string) (net.Conn, error) {
// 		if !trace.SpanFromContext(ctx).IsRecording() {
// 			return hook(ctx, network, addr)
// 		}
//
// 		ctx, span := th.conf.tracer.Start(ctx, "redis.dial", th.spanOpts...)
// 		defer span.End()
//
// 		conn, err := hook(ctx, network, addr)
// 		if err != nil {
// 			recordError(ctx, span, err)
// 			return nil, err
// 		}
// 		return conn, nil
// 	}
// }
//
// func (th *tracingHook) ProcessHook(hook redis.ProcessHook) redis.ProcessHook {
// 	return func(ctx context.Context, cmd redis.Cmder) error {
// 		if !trace.SpanFromContext(ctx).IsRecording() {
// 			return hook(ctx, cmd)
// 		}
//
// 		opts := th.spanOpts
// 		if th.conf.dbStmtEnabled {
// 			opts = append(opts, trace.WithAttributes(
// 				semconv.DBStatementKey.String(rediscmd.CmdString(cmd))),
// 			)
// 		}
//
// 		ctx, span := th.conf.tracer.Start(ctx, cmd.FullName(), opts...)
// 		defer span.End()
//
// 		if err := hook(ctx, cmd); err != nil {
// 			recordError(ctx, span, err)
// 			return err
// 		}
// 		return nil
// 	}
// }
//
// func (th *tracingHook) ProcessPipelineHook(
// 	hook redis.ProcessPipelineHook,
// ) redis.ProcessPipelineHook {
// 	return func(ctx context.Context, cmds []redis.Cmder) error {
// 		if !trace.SpanFromContext(ctx).IsRecording() {
// 			return hook(ctx, cmds)
// 		}
//
// 		opts := th.spanOpts
// 		opts = append(opts, trace.WithAttributes(
// 			attribute.Int("db.redis.num_cmd", len(cmds)),
// 		))
//
// 		summary, cmdsString := rediscmd.CmdsString(cmds)
// 		if th.conf.dbStmtEnabled {
// 			opts = append(opts, trace.WithAttributes(semconv.DBStatementKey.String(cmdsString)))
// 		}
//
// 		ctx, span := th.conf.tracer.Start(ctx, "redis.pipeline "+summary, opts...)
// 		defer span.End()
//
// 		if err := hook(ctx, cmds); err != nil {
// 			recordError(ctx, span, err)
// 			return err
// 		}
// 		return nil
// 	}
// }
//
// func recordError(ctx context.Context, span trace.Span, err error) {
// 	if err != redis.Nil {
// 		span.RecordError(err)
// 		span.SetStatus(codes.Error, err.Error())
// 	}
// }
//
// func formatDBConnString(network, addr string) string {
// 	if network == "tcp" {
// 		network = "redis"
// 	}
// 	return fmt.Sprintf("%s://%s", network, addr)
// }
