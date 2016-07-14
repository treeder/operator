package common

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

// WithLogger stores the logger in the ctx
func WithLogger(ctx context.Context, l logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, "logger", l)
}

// Logger returns the structured logger from the ctx
func Logger(ctx context.Context) logrus.FieldLogger {
	l, ok := ctx.Value("logger").(logrus.FieldLogger)
	if !ok {
		return logrus.StandardLogger()
	}
	return l
}

// Attempt at simplifying this whole logger in the context thing
// LoggerWithFields will return a logger with the new fields added and the logger will be set in the context
func LoggerWithFields(ctx context.Context, fields map[string]interface{}) (context.Context, logrus.FieldLogger) {
	l := Logger(ctx)
	l = l.WithFields(fields)
	ctx = WithLogger(ctx, l)
	return ctx, l
}

// Attempt at simplifying this whole logger in the context thing
func LoggerWithField(ctx context.Context, key, value string) (context.Context, logrus.FieldLogger) {
	l := Logger(ctx)
	l = l.WithField(key, value)
	ctx = WithLogger(ctx, l)
	return ctx, l
}
