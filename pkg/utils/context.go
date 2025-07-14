package utils

import (
	"context"
	"errors"
)

type userId string

const UserIDKey userId = "userID"

func GetUserIDFromContext(ctx context.Context) (string, error) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return "", errors.New("userID not found in context")
	}

	userID, ok := val.(string)
	if !ok {
		return "", errors.New("userID is of invalid type")
	}

	return userID, nil
}
