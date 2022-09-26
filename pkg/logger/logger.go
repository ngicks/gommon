package logger

import (
	"context"

	"go.uber.org/zap"
)

type ContextField struct {
	Label      string
	ContextKey interface{}
}

type Logger interface {
	Debug(i ...interface{})
	Debugf(format string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(i ...interface{})
	Infof(format string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(i ...interface{})
	Warnf(format string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(i ...interface{})
	Errorf(format string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatal(i ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Panic(i ...interface{})
	Panicf(format string, args ...interface{})
	Panicw(msg string, keysAndValues ...interface{})

	With(args ...interface{}) Logger
	// SetContextKey creates a new logger that inherites logging context and fields is overridden.
	// Fields are used when WithContext is called
	// to pull value from context.Context and to add label-value pairs to logging context.
	SetContextKey(fields ...ContextField) Logger
	// WithContext creates a child logger and adds strutured context.
	// It uses []ContextField which is set with SetContextKey to build logging context.
	// Structure is list of key-value pair
	// where key is ContextField.Label and value is retrieved by context#Value(ContextField.ContextKey).
	WithContext(ctx context.Context) Logger
}

var _ Logger = &ZapLogger{}

type ZapLogger struct {
	*zap.SugaredLogger
	contextFields []ContextField
}

func NewZapLogger(sugar *zap.SugaredLogger) *ZapLogger {
	return &ZapLogger{
		SugaredLogger: sugar,
	}
}

func (l *ZapLogger) WithRaw(args ...interface{}) *ZapLogger {
	return &ZapLogger{
		SugaredLogger: l.SugaredLogger.With(args...),
		contextFields: l.contextFields,
	}
}

func (l *ZapLogger) With(args ...interface{}) Logger {
	return l.WithRaw(args...)
}

func (l *ZapLogger) SetContextKeyRaw(fields ...ContextField) *ZapLogger {
	return &ZapLogger{
		SugaredLogger: l.SugaredLogger,
		contextFields: fields,
	}
}

func (l *ZapLogger) SetContextKey(fields ...ContextField) Logger {
	return l.SetContextKeyRaw(fields...)
}

func (l *ZapLogger) WithContextRaw(ctx context.Context) *ZapLogger {
	var labelAndValue []interface{}
	for _, v := range l.contextFields {
		labelAndValue = append(labelAndValue, v.Label, ctx.Value(v.ContextKey))
	}
	return l.WithRaw(labelAndValue...)
}

func (l *ZapLogger) WithContext(ctx context.Context) Logger {
	return l.WithContextRaw(ctx)
}
