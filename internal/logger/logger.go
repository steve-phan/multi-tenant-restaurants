package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

func Initialize(environment string) error {
	var config zap.Config
	if environment == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Logger, err = config.Build(zap.AddCaller())
	if err != nil {
		return err
	}
	Sugar = Logger.Sugar()
	return nil
}

func WithContext(ctx context.Context) *zap.Logger {
	fields := []zap.Field{}
	if rid := ctx.Value("request_id"); rid != nil {
		if id, ok := rid.(string); ok {
			fields = append(fields, zap.String("request_id", id))
		}
	}
	if uid := ctx.Value("user_id"); uid != nil {
		if id, ok := uid.(uint); ok {
			fields = append(fields, zap.Uint("user_id", id))
		}
	}
	return Logger.With(fields...)
}

func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}

func LogRequest(method, path string, status int, duration time.Duration) {
	Info("request",
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", status),
		zap.Duration("duration", duration),
	)
}

func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}
