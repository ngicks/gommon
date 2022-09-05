package logger

import "go.uber.org/zap"

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
}

var _ Logger = &ZapLogger{}

type ZapLogger struct {
	*zap.SugaredLogger
}

func NewZapLogger(sugar *zap.SugaredLogger) *ZapLogger {
	return &ZapLogger{
		SugaredLogger: sugar,
	}
}

func (l *ZapLogger) With(args ...interface{}) Logger {
	return NewZapLogger(l.SugaredLogger.With(args...))
}
