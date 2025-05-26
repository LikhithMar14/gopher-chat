package utils

import (
	"context"

)

type PostContextKey string
const PostIDKey PostContextKey = "post_id"

type UserContextKey string
const UserIDKey UserContextKey = "user_id"

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}


func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

