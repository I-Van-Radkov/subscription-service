package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	loggerRequestIDKey = "x-request-id"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Debug(ctx context.Context, msg string, fields ...zap.Field)
}

type L struct {
	z *zap.Logger
}

func NewLogger(env string) Logger {
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if env == "dev" {
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logger, err := loggerCfg.Build()
	if err != nil {
		return nil
	}
	defer logger.Sync()

	lo := L{logger}
	return &lo
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, loggerRequestIDKey, requestID)
}

func (l *L) Info(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)

	fields = append(fields, zap.String(loggerRequestIDKey, id))

	l.z.Info(msg, fields...)
}
func (l *L) Error(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)
	fields = append(fields, zap.String(loggerRequestIDKey, id))
	l.z.Error(msg, fields...)
}
func (l *L) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)
	fields = append(fields, zap.String(loggerRequestIDKey, id))
	l.z.Debug(msg, fields...)
}
