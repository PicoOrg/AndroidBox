package util

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg string, fields Fields)
	Debug(msg string, fields Fields)
	Warn(msg string, fields Fields)
	Error(msg string, fields Fields)
}

type logger struct {
	logger *zap.Logger
}

type Fields map[string]interface{}

func NewLogger(debug bool) Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if debug {
			return lvl >= zapcore.DebugLevel
		} else {
			return lvl >= zapcore.InfoLevel
		}
	})
	consoleLogger := zapcore.Lock(&WrapStdout{})
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, consoleLogger, level),
	}
	loggerInstance := zap.New(zapcore.NewTee(cores...), zap.AddCaller())
	return &logger{
		logger: loggerInstance,
	}
}

func (instance *logger) Info(msg string, fields Fields) {
	instance.log("info", msg, fields)
}

func (instance *logger) Debug(msg string, fields Fields) {
	instance.log("debug", msg, fields)
}

func (instance *logger) Warn(msg string, fields Fields) {
	instance.log("warn", msg, fields)
}

func (instance *logger) Error(msg string, fields Fields) {
	instance.log("error", msg, fields)
}

func (instance *logger) log(level, msg string, fields Fields) {
	zapFields := make([]zap.Field, 0)
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	switch level {
	case "debug":
		instance.logger.Debug(msg, zapFields...)
	case "info":
		instance.logger.Info(msg, zapFields...)
	case "warn":
		instance.logger.Warn(msg, zapFields...)
	case "error":
		instance.logger.Error(msg, zapFields...)
	}
}

type WrapStdout struct {
}

func (w *WrapStdout) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (w *WrapStdout) Sync() error {
	return os.Stdout.Sync()
}
