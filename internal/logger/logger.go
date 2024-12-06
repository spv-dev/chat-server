package logger

import (
	"flag"
	"log"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

var logLevel = flag.String("l", "info", "log level")

// Init инициализация логгера
func Init(core zapcore.Core, options ...zap.Option) {
	globalLogger = zap.New(core, options...)
}

// DefaultInit инициализация логгера с устновкой ядра
func DefaultInit(options ...zap.Option) {
	core := getCore(getAtomicLevel())
	globalLogger = zap.New(core, options...)
}

// Debug добавление в лог отладочной информации
func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

// Info добавление в лог информационного сообщения
func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

// Warn добавление в лог предупреждения
func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

// Error добавление в лог ошибки
func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

// Fatal добавление в лог фатальной ошибки
func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

// WithOptions добавление опций к логгеру
func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
		},
	)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
