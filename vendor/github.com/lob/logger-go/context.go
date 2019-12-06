package logger

import "context"

type key struct{}

// WithContext returns a copy of ctx with log attached to it.
func (log Logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, key{}, log)
}

// FromContext returns the Logger that is attached to ctx. If there is no
// Logger, a new Logger instance is returned.
func FromContext(ctx context.Context) Logger {
	if log, ok := ctx.Value(key{}).(Logger); ok {
		return log
	}
	return New()
}
