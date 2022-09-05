package echologgerzap

import (
	"io"
	"os"
	"sync"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapFactory = func(
	w io.Writer,
	logLevel zapcore.Level,
) *zap.Logger

type EchoLoggerZap struct {
	mu      sync.RWMutex
	w       io.Writer
	zap     *zap.SugaredLogger
	level   log.Lvl
	factory ZapFactory

	prefix string
}

// Default creates *EchoLoggerZap with opinionated defaults.
func Default() *EchoLoggerZap {
	var factory ZapFactory = func(
		w io.Writer,
		logLevel zapcore.Level,
	) *zap.Logger {
		conf := zap.NewProductionEncoderConfig()
		conf.EncodeTime = zapcore.RFC3339NanoTimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(conf),
			zapcore.AddSync(w),
			logLevel,
		)
		return zap.New(core, zap.AddCallerSkip(2), zap.AddCaller())
	}

	w := zapcore.AddSync(os.Stdout)

	return &EchoLoggerZap{
		w:       w,
		zap:     factory(w, zapcore.InfoLevel).Sugar(),
		level:   log.DEBUG,
		factory: factory,
	}
}

func New(w io.Writer, factory ZapFactory, level log.Lvl) *EchoLoggerZap {
	return &EchoLoggerZap{
		w:       w,
		factory: factory,
		level:   level,
		zap:     factory(w, zap.InfoLevel).Sugar(),
	}
}

func (l *EchoLoggerZap) Output() io.Writer {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.w
}
func (l *EchoLoggerZap) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.w = w
	l.zap = l.factory(l.w, echoToZapLevel(l.level)).Sugar()
}
func (l *EchoLoggerZap) Prefix() string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.prefix
}
func (l *EchoLoggerZap) SetPrefix(p string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.prefix = p
}
func (l *EchoLoggerZap) Level() log.Lvl {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.level
}
func (l *EchoLoggerZap) SetLevel(v log.Lvl) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.level = v
	l.zap = l.factory(l.w, echoToZapLevel(l.level)).Sugar()
}
func (l *EchoLoggerZap) SetHeader(h string) {
	l.mu.Lock()
	defer l.mu.Unlock()

}
func (l *EchoLoggerZap) Print(i ...interface{}) {
	l.log(zap.InfoLevel, i...)
}
func (l *EchoLoggerZap) Printf(format string, args ...interface{}) {
	l.logf(zap.InfoLevel, format, args...)
}
func (l *EchoLoggerZap) Printj(j log.JSON) {
	l.logj(zap.InfoLevel, j)
}
func (l *EchoLoggerZap) Debug(i ...interface{}) {
	l.log(zap.DebugLevel, i...)
}
func (l *EchoLoggerZap) Debugf(format string, args ...interface{}) {
	l.logf(zap.DebugLevel, format, args...)
}
func (l *EchoLoggerZap) Debugj(j log.JSON) {
	l.logj(zap.DebugLevel, j)
}
func (l *EchoLoggerZap) Info(i ...interface{}) {
	l.log(zap.InfoLevel, i...)
}
func (l *EchoLoggerZap) Infof(format string, args ...interface{}) {
	l.logf(zap.InfoLevel, format, args...)
}
func (l *EchoLoggerZap) Infoj(j log.JSON) {
	l.logj(zap.InfoLevel, j)
}
func (l *EchoLoggerZap) Warn(i ...interface{}) {
	l.log(zap.WarnLevel, i...)
}
func (l *EchoLoggerZap) Warnf(format string, args ...interface{}) {
	l.logf(zap.WarnLevel, format, args...)
}
func (l *EchoLoggerZap) Warnj(j log.JSON) {
	l.logj(zap.WarnLevel, j)
}
func (l *EchoLoggerZap) Error(i ...interface{}) {
	l.log(zap.ErrorLevel, i...)
}
func (l *EchoLoggerZap) Errorf(format string, args ...interface{}) {
	l.logf(zap.ErrorLevel, format, args...)
}
func (l *EchoLoggerZap) Errorj(j log.JSON) {
	l.logj(zap.ErrorLevel, j)
}
func (l *EchoLoggerZap) Fatal(i ...interface{}) {
	l.log(zap.FatalLevel, i...)
}
func (l *EchoLoggerZap) Fatalf(format string, args ...interface{}) {
	l.logf(zap.FatalLevel, format, args...)
}
func (l *EchoLoggerZap) Fatalj(j log.JSON) {
	l.logj(zap.FatalLevel, j)
}
func (l *EchoLoggerZap) Panic(i ...interface{}) {
	l.log(zap.PanicLevel, i...)
}
func (l *EchoLoggerZap) Panicf(format string, args ...interface{}) {
	l.logf(zap.PanicLevel, format, args...)
}
func (l *EchoLoggerZap) Panicj(j log.JSON) {
	l.logj(zap.PanicLevel, j)
}

func (l *EchoLoggerZap) log(level zapcore.Level, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var logger *zap.SugaredLogger
	if l.prefix != "" {
		logger = l.zap.With("prefix", l.prefix)
	} else {
		logger = l.zap
	}

	switch level {
	case zapcore.DebugLevel:
		logger.Debug(args...)
	case zapcore.InfoLevel:
		logger.Info(args...)
	case zapcore.WarnLevel:
		logger.Warn(args...)
	case zapcore.ErrorLevel:
		logger.Error(args...)
	case zapcore.DPanicLevel:
		logger.DPanic(args...)
	case zapcore.PanicLevel:
		logger.Panic(args...)
	case zapcore.FatalLevel:
		logger.Fatal(args...)
	}
}

func (l *EchoLoggerZap) logf(level zapcore.Level, format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var logger *zap.SugaredLogger
	if l.prefix != "" {
		logger = l.zap.With("prefix", l.prefix)
	} else {
		logger = l.zap
	}

	switch level {
	case zapcore.DebugLevel:
		logger.Debugf(format, args...)
	case zapcore.InfoLevel:
		logger.Infof(format, args...)
	case zapcore.WarnLevel:
		logger.Warnf(format, args...)
	case zapcore.ErrorLevel:
		logger.Errorf(format, args...)
	case zapcore.DPanicLevel:
		logger.DPanicf(format, args...)
	case zapcore.PanicLevel:
		logger.Panicf(format, args...)
	case zapcore.FatalLevel:
		logger.Fatalf(format, args...)
	}
}

func (l *EchoLoggerZap) logj(level zapcore.Level, j map[string]interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var logger *zap.SugaredLogger
	if l.prefix != "" {
		logger = l.zap.With("prefix", l.prefix)
	} else {
		logger = l.zap
	}

	var kv []interface{}
	for k, v := range j {
		kv = append(kv, k, v)
	}

	switch level {
	case zapcore.DebugLevel:
		logger.Debugw("", kv...)
	case zapcore.InfoLevel:
		logger.Infow("", kv...)
	case zapcore.WarnLevel:
		logger.Warnw("", kv...)
	case zapcore.ErrorLevel:
		logger.Errorw("", kv...)
	case zapcore.DPanicLevel:
		logger.DPanicw("", kv...)
	case zapcore.PanicLevel:
		logger.Panicw("", kv...)
	case zapcore.FatalLevel:
		logger.Fatalw("", kv...)
	}
}
