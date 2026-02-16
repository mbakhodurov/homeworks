package logger

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Key string

const (
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"
)

var (
	globalLogger *logger
	initOnce     sync.Once
	dynamicLevel zap.AtomicLevel
)

type logger struct {
	zapLogger *zap.Logger
}

func Init(levelStr string, asJSON bool) error {
	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(levelStr))

		encoderCfg := buildProductionEncoderConfig()

		var encoder zapcore.Encoder
		if asJSON {
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		}

		// 📁 Папка по дате
		dateFolder := time.Now().Format("2006-01-02")
		logDir := filepath.Join("logs", dateFolder)

		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			panic(err)
		}

		logFile := filepath.Join(logDir, "app.log")

		fileWriter := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    500, // 500 MB
			MaxBackups: 10,
			MaxAge:     30,   // дней хранения
			Compress:   true, // gzip старых логов
		}

		multiWriter := zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(fileWriter),
		)

		core := zapcore.NewCore(
			encoder,
			multiWriter,
			// zapcore.AddSync(os.Stdout),
			dynamicLevel,
		)

		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})
	return nil
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",                 // время
		LevelKey:       "level",                     // уровень логирования
		NameKey:        "logger",                    // имя логгера, если используется
		CallerKey:      "caller",                    // откуда вызван лог
		MessageKey:     "message",                   // текст сообщения
		StacktraceKey:  "stacktrace",                // стектрейс для ошибок
		LineEnding:     zapcore.DefaultLineEnding,   // перенос строки
		EncodeLevel:    zapcore.CapitalLevelEncoder, // INFO, ERROR
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // читаемый ISO 8601 формат
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // короткий caller
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func SetLevel(levelStr string) {
	if dynamicLevel == (zap.AtomicLevel{}) {
		return
	}
	dynamicLevel.SetLevel(parseLevel(levelStr))
}

func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func Logger() *logger {
	return globalLogger
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Debug(ctx, msg, fields...)
}
func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allField := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allField...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Error(ctx, msg, fields...)
}
func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Info(ctx, msg, fields...)
}
func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Warn(ctx, msg, fields...)
}
func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Fatal(ctx, msg, fields...)
}
func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}

// fieldsFromContext вытаскивает enrich-поля из контекста
func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}
