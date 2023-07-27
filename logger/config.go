package logger

import ()

type OutputQueryLogType string

const (
	// output log nothing
	OutputQueryLogTypeNone OutputQueryLogType = "none"

	// output log all query
	OutputQueryLogTypeAll OutputQueryLogType = "all"

	// output log only for select query
	OutputQueryLogTypeReadOnly OutputQueryLogType = "read"

	// output log only for insert/update/delete query
	OutputQueryLogTypeWriteOnly OutputQueryLogType = "write"
)

type config struct {
	fields             map[string]any
	outputQueryLogType OutputQueryLogType
}

type Option interface {
	apply(conf *config)
}

type option func(conf *config)

func (fn option) apply(conf *config) {
	fn(conf)
}

func newConfig(opts ...Option) *config {
	conf := &config{
		outputQueryLogType: OutputQueryLogTypeAll,
	}

	for _, opt := range opts {
		opt.apply(conf)
	}

	return conf
}

func WithLogFields(fields map[string]any) Option {
	return option(func(conf *config) {
		conf.fields = fields
	})
}

func WithOutputQueryLogType(t OutputQueryLogType) Option {
	return option(func(conf *config) {
		conf.outputQueryLogType = t
	})
}

func (t OutputQueryLogType) AllowReadLog() bool {
	return t == OutputQueryLogTypeAll || t == OutputQueryLogTypeReadOnly
}

func (t OutputQueryLogType) AllowWriteLog() bool {
	return t == OutputQueryLogTypeAll || t == OutputQueryLogTypeWriteOnly
}
