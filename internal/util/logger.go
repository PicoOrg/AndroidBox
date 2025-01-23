package util

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel int8 = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
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

func NewLogger(level int8) Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleLogger := zapcore.Lock(&WrapStdout{})
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, consoleLogger, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return int8(lvl) >= level
		})),
	}
	loggerInstance := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(2))
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
