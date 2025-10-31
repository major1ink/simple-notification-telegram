package logger

import (
	"time"

	"go.uber.org/zap/zapcore"
)

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func getLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case "debug":
		return zapcore.DebugLevel
	case "DEBUG":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "INFO":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "WARN":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
