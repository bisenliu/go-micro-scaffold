package contextutil

import (
	"context"
)

// UserIDKey is the key for the user ID in the context.
const UserIDKey = "user_id"

// GetUserIDFromContext retrieves the user ID from the context.
// It returns the user ID and a boolean indicating whether the user ID was found.
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
