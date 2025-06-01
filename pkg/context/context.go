package context

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

func SetLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

func GetLogger(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(LoggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return zap.NewNop().Sugar()
}

func GetLoggerWithDefaults(ctx context.Context, defaultLogger *zap.SugaredLogger) *zap.SugaredLogger {
	if logger, ok := ctx.Value(LoggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return defaultLogger
}
