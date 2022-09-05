package echologgerzap

import (
	"github.com/labstack/gommon/log"
	"go.uber.org/zap/zapcore"
)

func echoToZapLevel(lvl log.Lvl) zapcore.Level {
	switch lvl {
	case log.DEBUG:
		return zapcore.DebugLevel
	case log.INFO:
		return zapcore.InfoLevel
	case log.WARN:
		return zapcore.WarnLevel
	case log.ERROR:
		return zapcore.ErrorLevel
	case log.OFF:
		return zapcore.FatalLevel
	default:
		return zapcore.FatalLevel
		// case 6: //log.panicLevel
		// case 7: //log.fatalLevel
	}
}
