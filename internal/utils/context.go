package utils

import (
	"context"

	"go.uber.org/zap"
)

type PostContextKey string

const PostIDKey PostContextKey = "post_id"

type UserContextKey string

const UserIDKey UserContextKey = "user_id"

type LoggerContextKey string

const LoggerKey LoggerContextKey = "logger"

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

// SetLogger adds a logger to the context
func SetLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// GetLogger retrieves the logger from context
func GetLogger(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(LoggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	// Return a no-op logger if none found (fallback)
	return zap.NewNop().Sugar()
}

// GetLoggerWithDefaults retrieves logger from context or returns a default one
func GetLoggerWithDefaults(ctx context.Context, defaultLogger *zap.SugaredLogger) *zap.SugaredLogger {
	if logger, ok := ctx.Value(LoggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return defaultLogger
}
